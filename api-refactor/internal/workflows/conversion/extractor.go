package conversion

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
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
	l.mu.Lock() // calling routine blocks other routines from modifying the mutex
	defer l.mu.Unlock()
	l.e.SuccessFiles[f] = b
} // logSuccess

func (l *Locker) logFailure(f shared.FileID, msg string) {
	l.mu.Lock() // calling routine blocks other routines from modifying the mutex
	defer l.mu.Unlock()
	l.e.FailedFiles[f] = msg

} // logFailure

func (c *Conversion) Extract(ctx context.Context, d *shared.DownloadResult) (*ExtractionResult, error) {
	var err error
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
	log.Printf("Go path: %s", goDir)
	// build path to python venv and script
	pythonDir := filepath.Join(goDir, "python")
	venvDir := filepath.Join(pythonDir, "venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "convert_pdf.py")
	log.Printf("Python script path: %s", scriptPath)

	cmd := exec.Command(pythonExec, "-u", scriptPath)
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

		for {
			metaLine, err := bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				} // if
				log.Printf("read stdout meta line error: %v", err)
				return
			} // if
			metaLine = metaLine[:len(metaLine)-1] // strips the newline character

			var meta FileJSON
			if err := json.Unmarshal([]byte(metaLine), &meta); err != nil {
				log.Printf("invalid metadata JSON from python %v -- raw: %s", err, metaLine)
				continue
			} // if

			// read b64 body until sentinel
			var buf bytes.Buffer
			sentinel := []byte(bodySentinel)
			tmp := make([]byte, 4096)
			for {
				n, readErr := bufReader.Read(tmp)
				if n > 0 {
					chunk := tmp[:n]
					// check if sentinel is in chunk
					if idx := bytes.Index(chunk, sentinel); idx >= 0 {
						// read up until sentinel and exit loop
						buf.Write(chunk[:idx])
						break
					} // if
					// otherwise, write up to buf size
					buf.Write(chunk)
				} // if
				if readErr != nil {
					if readErr == io.EOF {
						break
					} // if
					log.Printf("error reading stdout body: %v", readErr)
					break
				} // if
			} // for

			// buf now contains b64 encoded raw text
			decoded, decErr := base64.StdEncoding.DecodeString(buf.String())
			if decErr != nil {
				msg := fmt.Sprintf("unable to decode base64 body from python: %v", decErr)
				log.Print(msg)
				fid := shared.FileID{ObjectKey: meta.ObjectKey}
				l.logFailure(fid, msg)
				continue
			} // if

			fileToAdd := shared.FileID{
				Title:       meta.Title,
				ObjectKey:   meta.ObjectKey,
				OriginalKey: meta.OriginalKey,
				URL:         meta.URL,
			}
			l.logSuccess(fileToAdd, decoded)
		} // for

		//	// Reading converted files line by line
		//	scanner := bufio.NewScanner(stdout)
		//	for scanner.Scan() {
		//		var success FileJSON
		//		var fileToAdd shared.FileID
		//		if err := json.Unmarshal(scanner.Bytes(), &success); err != nil {
		//			log.Printf("error decoding response: %v", err)
		//			if e, ok := err.(*json.SyntaxError); ok {
		//				log.Printf("syntax error at byte offset %d, %v", e.Offset, success)
		//			} // if
		//		} // if

		//		fileToAdd = shared.FileID{
		//			Title:       success.Title,
		//			ObjectKey:   success.ObjectKey,
		//			OriginalKey: success.OriginalKey,
		//			URL:         success.URL,
		//		}
		//		utf8body := string(success.Body)
		//		decodedBody, err := base64.StdEncoding.DecodeString(utf8body)
		//		if err != nil {
		//			log.Panic("unable to decode properly processed file")
		//		} // if
		//		l.logSuccess(fileToAdd, decodedBody)
		//	} // for

	} // readStdout

	go readStdout(&l)

	readStderr := func(l *Locker) {
		defer wg.Done()
		stderrScan := bufio.NewScanner(stderr)

		for stderrScan.Scan() {
			line := stderrScan.Text()
			// try to parse
			var failure FileJSON
			if err := json.Unmarshal([]byte(line), &failure); err == nil && failure.Err != "" {
				fid := shared.FileID{ObjectKey: failure.Err}
				l.logFailure(fid, failure.Err)
			} else {
				log.Printf("python stderr: %s", line)
			} // if
		} // for
		if err := stderrScan.Err(); err != nil {
			log.Printf("stderr scanner error: %v", err)
		} // if

		//// Reading failed files
		//scanner := bufio.NewScanner(stderr)
		//for scanner.Scan() {
		//	var failure FileJSON
		//	var fileToAdd shared.FileID
		//	// correctly read error
		//	if json.Valid(scanner.Bytes()) {
		//		if err := json.Unmarshal(scanner.Bytes(), &failure); err != nil {
		//			fileToAdd = shared.FileID{
		//				ObjectKey: failure.Err,
		//			}
		//			// log the file error
		//			l.logFailure(fileToAdd, fileToAdd.ObjectKey)
		//		} // if
		//	} else {
		//		// NOTE: THIS SHOULD NEVER THROW
		//		log.Fatalf("Non-JSON error: %s", scanner.Text())
		//		if err := json.Unmarshal(scanner.Bytes(), &failure); err != nil {
		//			if e, ok := err.(*json.SyntaxError); ok {
		//				log.Printf("syntax error at byte offset %d, %v", e.Offset, failure)
		//			} // if
		//		} // if
		//	} // if
		//} // for

	} // readStderr

	go readStderr(&l)

	// write json objects
	for id, r := range files {
		// build metadata
		jsonObject := FileJSON{
			Body:        "",
			Title:       id.Title,
			ObjectKey:   id.ObjectKey,
			URL:         id.URL,
			OriginalKey: id.ObjectKey,
		}

		// build metadata
		metaJSON, err := json.Marshal(jsonObject)
		if err != nil {
			msg := fmt.Sprintf("marshal metadata: %v", err)
			l.logFailure(id, msg)
		} // if

		// write metadata + newline
		if _, err := stdin.Write(metaJSON); err != nil {
			msg := fmt.Errorf("write meta to stdin: %v", err)
			return nil, msg
		} // if
		if _, err := stdin.Write([]byte("\n")); err != nil {
			msg := fmt.Errorf("write newline to stdin: %v", err)
			return nil, msg
		} // if

		// wrap stdin with b64 encoder
		encoder := base64.NewEncoder(base64.StdEncoding, stdin)
		// stream file to python subprocess
		_, err = io.Copy(encoder, r)

		// close encoder so b64 padding is applied and input is flushed
		if cerr := encoder.Close(); cerr != nil {
			log.Printf("encoder close error: %v", cerr)
		} // if

		// check error on io copy
		if err != nil {
			msg := fmt.Sprintf("io copy to python subprocess: %v", err)
			l.logFailure(id, msg)
			continue
		} // if

		// write file delimeter
		if _, err := stdin.Write([]byte(bodySentinel)); err != nil {
			msg := fmt.Errorf("write sentinel: %v", err)
			return nil, msg
		} // if

		// close reader
		if err := r.Close(); err != nil {
			msg := fmt.Sprintf("close reader: %v", err)
			l.logFailure(id, msg)
		} // if
	} // for

	// close stdin after all files are written
	if err := stdin.Close(); err != nil {
		msg := fmt.Errorf("stdin close error: %v", err)
		return nil, msg
	} // if

	// wait for readers
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		log.Printf("python exit error: %v", err)
	} // if
	return &l.e, nil
} // Extract

/*
Returns a properly encoded json object ready for streaming

On error, populates ExtractionResult.FailedFiles
*/
//func buildJSON(id shared.FileID, r io.ReadCloser) ([]byte, error) {
//
//	var err error
//	var buf []byte
//	var jsonData []byte
//	var jsonObject FileJSON
//	buf, err = io.ReadAll(r)
//	if err != nil {
//		return jsonData, err
//	} // if
//
//	// encode body in b64
//	encodedBuf := base64.StdEncoding.EncodeToString(buf)
//
//	// encode data
//	jsonObject = FileJSON{
//		Body:        encodedBuf,
//		Title:       id.Title,
//		ObjectKey:   id.ObjectKey,
//		URL:         id.URL,
//		OriginalKey: id.ObjectKey,
//	}
//	jsonData, err = json.Marshal(jsonObject)
//	if err != nil {
//		return jsonData, err
//	} // if
//	return jsonData, err
//} // buildJSON
