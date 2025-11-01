package conversion

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
	"sync"

	"tui/backend/types/shared"
)

type ExtractionResult struct {
	FailedFiles  map[shared.FileID]string
	SuccessFiles map[shared.FileID][]byte
}

type FileJSON struct {
	Body      string `json:"body,omitempty"`
	Title     string `json:"title,omitempty"`
	ObjectKey string `json:"objectKey,omitempty"`
	URL       string `json:"url,omitempty"`
	Err       string `json:"error,omitempty"`
}

type Locker struct {
	mu sync.Mutex
	e  ExtractionResult
}

/*
NOTE:

I'm not sure we'll be able to use an any type here.
The encoding/decoding process on the python side of things
might behave weirdly.
*/
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
	projectRoot, _ := os.Getwd()
	pythonDir := filepath.Join(projectRoot, "python")
	venvDir := filepath.Join(pythonDir, "venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "convert_pdf.py")
	cmd := exec.Command(pythonExec, "-u", scriptPath)
	cmd.Dir = pythonDir

	// copy curr environment, but inject venv info
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
		fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	)
	// wait group for routines
	var wg sync.WaitGroup
	// pipes for processing
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// start command execution
	if err := cmd.Start(); err != nil {
		panic(err)
	} // if

	// tell the waitgroup we're waiting for two routines to finish
	wg.Add(2)

	readStdout := func(l *Locker) {
		defer wg.Done()

		// Reading converted files line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var success FileJSON
			var fileToAdd shared.FileID
			if err := json.Unmarshal(scanner.Bytes(), &success); err != nil {
				log.Printf("error decoding response: %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("syntax error at byte offset %d, %v", e.Offset, success)
				} // if
			} // if

			fileToAdd = shared.FileID{
				Title:     success.Title,
				ObjectKey: success.ObjectKey,
				URL:       success.URL,
			}
			utf8body := string(success.Body)
			decodedBody, err := base64.StdEncoding.DecodeString(utf8body)
			if err != nil {
				log.Panic("unable to decode properly processed file")
			} // if
			l.logSuccess(fileToAdd, decodedBody)
		} // for
	} // readStdout

	go readStdout(&l)

	readStderr := func(l *Locker) {
		defer wg.Done()

		// Reading failed files
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			var failure FileJSON
			var fileToAdd shared.FileID
			// correctly read error
			if json.Valid(scanner.Bytes()) {
				if err := json.Unmarshal(scanner.Bytes(), &failure); err != nil {
					fileToAdd = shared.FileID{
						ObjectKey: failure.Err,
					}
					// log the file error
					l.logFailure(fileToAdd, fileToAdd.ObjectKey)
				} // if
			} else {
				// NOTE: THIS SHOULD NEVER THROW
				log.Fatalf("Non-JSON error: %s", scanner.Text())
				if err := json.Unmarshal(scanner.Bytes(), &failure); err != nil {
					if e, ok := err.(*json.SyntaxError); ok {
						log.Printf("syntax error at byte offset %d, %v", e.Offset, failure)
					} // if
				} // if
			} // if
		} // for
	} // readStderr

	go readStderr(&l)

	// write json objects
	for id, r := range files {
		jsonData, err := buildJSON(id, r)

		// was not able to build json from file info, log and continue
		if err != nil {
			msg := fmt.Sprintf("could not build json from %s", id.Title)
			l.logFailure(id, msg)
			continue
		} // if
		if _, err := stdin.Write(jsonData); err != nil {
			// we'll panic here because there is something wrong
			// with json object if we can't write it
			panic(err)
		} // if
		if _, err = stdin.Write([]byte("\n")); err != nil {
			panic(err)
		} // if
	} // for
	if err := stdin.Close(); err != nil {
		panic(err)
	} // if

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		panic(err)
	} // if
	return &l.e, err
} // Extract

/*
Returns a properly encoded json object ready for streaming

On error, populates ExtractionResult.FailedFiles
*/
func buildJSON(id shared.FileID, r io.ReadCloser) ([]byte, error) {

	var err error
	var buf []byte
	var jsonData []byte
	var jsonObject FileJSON
	buf, err = io.ReadAll(r)
	if err != nil {
		return jsonData, err
	} // if

	// encode body in b64
	encodedBuf := base64.StdEncoding.EncodeToString(buf)

	// encode data
	jsonObject = FileJSON{
		Body:      encodedBuf,
		Title:     id.Title,
		ObjectKey: id.ObjectKey,
		URL:       id.URL,
		//Err:       "",
	}
	jsonData, err = json.Marshal(jsonObject)
	if err != nil {
		return jsonData, err
	} // if
	return jsonData, err
} // buildJSON
