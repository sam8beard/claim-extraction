package processing

import (
	"bufio"
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
	"strings"
	"sync"
)

const bodySentinel = "--END-BODY--\n"

type PythonInput struct {
	ObjectKey string `json:"objectKey,omitempty"`
	Error     string `json:"error,omitempty"`
	FileName  string `json:"fileName,omitempty"`
}

type FileData struct {
	ObjectKey  string      `json:"objectKey,omitempty"`
	FileName   string      `json:"fileName,omitempty"`
	ClaimScore float64     `json:"claimScore,omitempty"`
	ClaimSpans []ClaimSpan `json:"claimSpans,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type ClaimSpan struct {
	Text       string  `json:"text"`
	Type       string  `json:"type"`
	Sent       string  `json:"sent"`
	Confidence float64 `json:"confidence"`
}

type NLPResult struct {
	FileData []FileData `json:"files"`
	mu       sync.Mutex
}

func (n *NLPResult) addFileData(fd FileData) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.FileData = append(n.FileData, fd)
}

func (p *Processing) NLP(ctx context.Context, f *FetchResult) (*NLPResult, error) {
	result := &NLPResult{
		FileData: make([]FileData, 0),
	}

	files := f.SuccessFiles

	/*
		Manually set python execution path and script path. By default,
		exec.Command uses the system environment for the executable.
		We need to use the environment located in venv/bin/python3 in
		order to execute the script so the libraries installed inside
		the virtual environment can be recognized.
	*/
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	pythonDir := filepath.Join(currentDir, "python")
	venvDir := filepath.Join(pythonDir, ".venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "nlp_processing.py")

	cmd := exec.Command(pythonExec, "-u", "-W", "ignore", scriptPath)
	cmd.Dir = pythonDir

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
		fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	)

	// Initialize pipes for processing
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

	log.Println("Executing NLP script...")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start script error: %w", err)
	}

	// Add routines to waitgroup
	var wg sync.WaitGroup
	wg.Add(2)

	readStdout := func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)

		// Increase buffer size for potentially large JSON lines
		buf := make([]byte, 4096*4096)
		scanner.Buffer(buf, 4096*4096)

		for scanner.Scan() {
			var fd FileData
			if err := json.Unmarshal(scanner.Bytes(), &fd); err != nil {
				log.Printf("bad JSON from python stdout: %v", err)
				continue
			}
			result.addFileData(fd)
			log.Printf("NLP processed file successfully: %s", fd.ObjectKey)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("stdout scanner error: %v", err)
		}
	}

	log.Println("Begin reading from stdout...")
	go readStdout()

	// Read stderr - log errors
	readStderr := func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			line := scanner.Text()
			var fd FileData

			if err := json.Unmarshal([]byte(line), &fd); err == nil && fd.Error != "" {
				log.Printf("NLP processing error for file %s: %s", fd.ObjectKey, fd.Error)
			} else {
				log.Printf("python stderr: %s", line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("stderr scanner error: %v", err)
		}
	}

	log.Println("Begin reading from stderr...")
	go readStderr()

	// Use buffered writer for stdin
	stdinWriter := bufio.NewWriter(stdin)

	// Write files to Python subprocess
	for file, reader := range files {
		// Get file name from object key
		fileName := file.ObjectKey
		fileName = strings.Replace(fileName, "processed/", "", 1)

		jsonObject := PythonInput{
			FileName:  fileName,
			ObjectKey: file.ObjectKey,
			Error:     "",
		}

		// Marshal metadata
		metaJSON, err := json.Marshal(jsonObject)
		if err != nil {
			log.Printf("marshal error for %s: %v", fileName, err)
			continue
		}

		// Write metadata + newline
		log.Printf("writing metadata for %s...", fileName)
		if _, err := stdinWriter.Write(metaJSON); err != nil {
			log.Printf("write meta to stdin: %v", err)
			continue
		}
		if _, err := stdinWriter.Write([]byte("\n")); err != nil {
			log.Printf("write newline to stdin: %v", err)
			continue
		}

		// Flush metadata before body
		if err := stdinWriter.Flush(); err != nil {
			log.Printf("flush metadata: %v", err)
			continue
		}

		// Write base64 encoded body directly to stdin
		encoder := base64.NewEncoder(base64.StdEncoding, stdin)

		log.Printf("writing file %s", fileName)
		_, err = io.Copy(encoder, reader)

		// Close encoder
		if cerr := encoder.Close(); cerr != nil {
			log.Printf("encoder close error: %v", cerr)
			if err := reader.Close(); err != nil {
				log.Printf("close reader after encoder error: %v", err)
			}
			continue
		}

		// Check error on copy
		if err != nil {
			log.Printf("IO copy to python subprocess: %v", err)
			if err := reader.Close(); err != nil {
				log.Printf("close reader after copy error: %v", err)
			}
			continue
		}

		// Write file delimiter
		if _, err := stdin.Write([]byte(bodySentinel)); err != nil {
			log.Printf("write sentinel error: %v", err)
			if err := reader.Close(); err != nil {
				log.Printf("close reader after sentinel error: %v", err)
			}
			continue
		}

		// Close reader
		if err := reader.Close(); err != nil {
			log.Printf("close reader: %v", err)
		}
	}

	log.Println("Closing stdin...")
	if err := stdin.Close(); err != nil {
		return nil, fmt.Errorf("stdin close error: %w", err)
	}

	// Wait for goroutines to finish
	log.Println("Waiting for goroutines...")
	wg.Wait()

	// Wait for script to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("python exit error: %v", err)
	}

	log.Printf("NLP processing complete. Processed %d files", len(result.FileData))
	return result, nil
}
