package conversion

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/sam8beard/claim-extraction/internal/types/shared"
)

const bodySentinel = "--END-BODY--\n"

type ExtractionResult struct {
	FailedFiles  map[shared.FileID]string
	SuccessFiles map[shared.FileID][]byte
}

type FileJSON struct {
	Body        string `json:"body,omitempty"`
	Title       string `json:"title,omitempty"`
	ObjectKey   string `json:"objectKey,omitempty"`
	URL         string `json:"url,omitempty"`
	Err         string `json:"error,omitempty"`
	OriginalKey string `json:"originalKey,omitempty"`
}

type Locker struct {
	mu sync.Mutex
	e  ExtractionResult
}

func (l *Locker) logSuccess(f shared.FileID, b []byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.e.SuccessFiles[f] = b
}

func (l *Locker) logFailure(f shared.FileID, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.e.FailedFiles[f] = msg
}

func (c *Conversion) Extract(ctx context.Context, d *shared.DownloadResult) (*ExtractionResult, error) {
	l := Locker{
		e: ExtractionResult{
			FailedFiles:  make(map[shared.FileID]string),
			SuccessFiles: make(map[shared.FileID][]byte),
		},
	}

	files := d.SuccessFiles

	/*
		Manually set python execution path and script path. By default,
		exec.Command uses the system environment for the executable.
		We need to use the environment located in venv/bin/python3 in
		order to execute the script so the libraries installed inside
		the virtual environment can be recognized.
	*/

	// get dir of curr go file
	_, goFile, _, _ := runtime.Caller(0)
	goDir := filepath.Dir(goFile)
	// build path to python venv and script
	pythonDir := filepath.Join(goDir, "python")
	venvDir := filepath.Join(pythonDir, ".venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "convert_pdf.py")

	cmd := exec.CommandContext(ctx, pythonExec, "-u", "-W", "ignore", scriptPath)
	cmd.Dir = pythonDir

	// copy curr environment, but inject venv info
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
		fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	)

	// initialize pipes for processing
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	// start command execution
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start python: %w", err)
	}

	// wait group for routines
	var wg sync.WaitGroup
	wg.Add(2)

	readStdout := func(l *Locker) {
		defer wg.Done()

		bufReader := bufio.NewReader(stdout)
		leftover := []byte{}

		for {
			// Read metadata line (handling leftover from previous iteration)
			var metaLine []byte

			if len(leftover) > 0 {
				// Check if leftover contains a complete line
				if idx := bytes.IndexByte(leftover, '\n'); idx >= 0 {
					metaLine = leftover[:idx]
					leftover = leftover[idx+1:]
				} else {
					// Need more data - read and append to leftover
					chunk, readErr := bufReader.ReadBytes('\n')
					if readErr != nil {
						if readErr == io.EOF {
							return
						}
						log.Printf("read stdout meta line error: %v", readErr)
						return
					}
					metaLine = append(leftover, bytes.TrimSuffix(chunk, []byte("\n"))...)
					leftover = nil
				}
			} else {
				// Read fresh line
				line, readErr := bufReader.ReadBytes('\n')
				if readErr != nil {
					if readErr == io.EOF {
						return
					}
					log.Printf("read stdout meta line error: %v", readErr)
					return
				}
				metaLine = bytes.TrimSpace(line)
			}

			if len(metaLine) == 0 {
				continue
			}

			// Parse metadata JSON
			var meta FileJSON
			if err := json.Unmarshal(metaLine, &meta); err != nil {
				log.Printf("invalid metadata JSON from python: %v", err)
				continue
			}

			// Read base64 body until sentinel
			var buf bytes.Buffer
			overlap := leftover
			leftover = nil

			for {
				chunk := make([]byte, 4096)
				n, err := bufReader.Read(chunk)
				if n > 0 {
					chunk = chunk[:n]
					combined := append(overlap, chunk...)

					// Check for sentinel
					idx := bytes.Index(combined, []byte(bodySentinel))
					if idx >= 0 {
						// Found sentinel - write everything before it
						buf.Write(combined[:idx])

						// Save everything after sentinel as leftover for next file
						start := idx + len(bodySentinel)
						if start < len(combined) {
							leftover = combined[start:]
						}
						break
					}

					// No sentinel found - use sliding window to handle sentinel split across chunks
					if len(combined) >= len(bodySentinel)-1 {
						// Write all but the last (sentinelLen-1) bytes
						writeUntil := len(combined) - (len(bodySentinel) - 1)
						buf.Write(combined[:writeUntil])
						overlap = combined[writeUntil:]
					} else {
						// Combined is smaller than sentinel - keep everything
						overlap = combined
					}
				}

				if err != nil {
					if err == io.EOF {
						log.Printf("unexpected EOF while reading body")
						break
					}
					log.Printf("error reading stdout body: %v", err)
					break
				}
			}

			// Decode base64 body
			decoded, decErr := base64.StdEncoding.DecodeString(buf.String())
			if decErr != nil {
				msg := fmt.Sprintf("unable to decode base64 body from python: %v", decErr)
				log.Print(msg)
				fid := shared.FileID{ObjectKey: meta.ObjectKey}
				l.logFailure(fid, msg)
				continue
			}

			fileToAdd := shared.FileID{
				Title:       meta.Title,
				ObjectKey:   meta.ObjectKey,
				OriginalKey: meta.OriginalKey,
				URL:         meta.URL,
			}
			l.logSuccess(fileToAdd, decoded)
		}
	}

	go readStdout(&l)

	// readStderr processes error messages from Python
	readStderr := func(l *Locker) {
		defer wg.Done()
		stderrScan := bufio.NewScanner(stderr)

		for stderrScan.Scan() {
			line := stderrScan.Text()
			// Try to parse as JSON error
			var failure FileJSON
			if err := json.Unmarshal([]byte(line), &failure); err == nil && failure.Err != "" {
				// Use OriginalKey if available, otherwise use ObjectKey from error
				objKey := failure.OriginalKey
				if objKey == "" {
					objKey = failure.ObjectKey
				}
				fid := shared.FileID{ObjectKey: objKey}
				l.logFailure(fid, failure.Err)
			} else {
				log.Printf("python stderr: %s", line)
			}
		}
		if err := stderrScan.Err(); err != nil {
			log.Printf("stderr scanner error: %v", err)
		}
	}

	go readStderr(&l)

	// Use buffered writer for stdin
	stdinWriter := bufio.NewWriter(stdin)

	// Write JSON objects to Python subprocess
	for id, r := range files {
		// Build metadata
		jsonObject := FileJSON{
			Body:        "",
			Title:       id.Title,
			ObjectKey:   id.ObjectKey,
			URL:         id.URL,
			OriginalKey: id.ObjectKey,
		}

		// Marshal metadata
		metaJSON, err := json.Marshal(jsonObject)
		if err != nil {
			msg := fmt.Sprintf("marshal metadata: %v", err)
			l.logFailure(id, msg)
			continue
		}

		// Write metadata + newline
		if _, err := stdinWriter.Write(metaJSON); err != nil {
			msg := fmt.Sprintf("write meta to stdin: %v", err)
			l.logFailure(id, msg)
			continue
		}
		if _, err := stdinWriter.Write([]byte("\n")); err != nil {
			msg := fmt.Sprintf("write newline to stdin: %v", err)
			l.logFailure(id, msg)
			continue
		}

		// Flush metadata before writing body
		if err := stdinWriter.Flush(); err != nil {
			msg := fmt.Sprintf("flush metadata: %v", err)
			l.logFailure(id, msg)
			continue
		}

		// Write base64 encoded body directly to stdin (not buffered writer)
		encoder := base64.NewEncoder(base64.StdEncoding, stdin)

		// Stream file to python subprocess
		_, err = io.Copy(encoder, r)

		// Close encoder so base64 padding is applied and input is flushed
		if cerr := encoder.Close(); cerr != nil {
			msg := fmt.Sprintf("encoder close error: %v", cerr)
			l.logFailure(id, msg)
			if err := r.Close(); err != nil {
				log.Printf("close reader after encoder error: %v", err)
			}
			continue
		}

		// Check error on io.Copy
		if err != nil {
			msg := fmt.Sprintf("io copy to python subprocess: %v", err)
			l.logFailure(id, msg)
			if err := r.Close(); err != nil {
				log.Printf("close reader after copy error: %v", err)
			}
			continue
		}

		// Write file delimiter
		if _, err := stdin.Write([]byte(bodySentinel)); err != nil {
			msg := fmt.Errorf("write sentinel: %w", err)
			return nil, msg
		}

		// Close reader
		if err := r.Close(); err != nil {
			msg := fmt.Sprintf("close reader: %v", err)
			l.logFailure(id, msg)
		}
	}

	// Close stdin after all files are written
	if err := stdin.Close(); err != nil {
		msg := fmt.Errorf("stdin close error: %w", err)
		return nil, msg
	}

	// Wait for goroutines to finish
	wg.Wait()

	// Wait for python process to exit
	if err := cmd.Wait(); err != nil {
		log.Printf("python exit error: %v", err)
	}

	// Log all failures
	for fid, errMsg := range l.e.FailedFiles {
		log.Printf("file failed to convert [%s]: %s", fid.ObjectKey, errMsg)
	}

	return &l.e, nil
}
