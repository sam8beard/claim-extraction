package conversion

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"tui/backend/types/shared"
)

type ExtractionResult struct {
	FailedFiles  map[shared.FileID]string
	SuccessFiles map[shared.FileID][]byte
}

// type entry {

// }

type FileJSON struct {
	Body      string `json:"body"`
	Title     string `json:"title"`
	ObjectKey string `json:"objectKey"`
	URL       string `json:"url"`
	Err       string `json:"error"`
}

type Locker struct {
	mu sync.Mutex
	e  ExtractionResult
}

func (l *Locker) log(f shared.FileID, u any) {
	l.mu.Lock() // calling routine blocks other routines from modifying the mutex
	defer l.mu.Unlock()

	switch val := u.(type) {
	case string:
		// this is a failed file
		l.e.FailedFiles[f] = val
	case []byte:
		// this is a successful file
		l.e.SuccessFiles[f] = val
	} // switch

	// if unknown
	// if  == empty {
	// 	// this is an unsuccessful file

	// } else {
	// 	// this is a successful file
	// 	fileToAdd := shared.FileID{
	// 		Title:     string(success.Title),
	// 		ObjectKey: string(success.ObjectKey),
	// 		URL:       string(success.URL),
	// 	}
	// 	// add successfully converted file
	// 	e.SuccessFiles[fileToAdd] = []byte(success.Body)
	// }
}

func (c *Conversion) Extract(ctx context.Context, d *shared.DownloadResult) (*ExtractionResult, error) {
	var err error

	l := Locker{
		e: ExtractionResult{
			FailedFiles:  make(map[shared.FileID]string),
			SuccessFiles: make(map[shared.FileID][]byte),
		},
	}

	var wg sync.WaitGroup

	files := d.SuccessFiles
	cmd := exec.Command("python3", "/python/convert_pdf.py")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	} // if
	// tell the waitgroup we're waiting for two routines to finish
	wg.Add(2)

	readStdout := func() {
		wg.Done()
		// Reading converted files line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var success FileJSON
			if err := json.Unmarshal(scanner.Bytes(), &success); err != nil {
				// THIS SHOULD NEVER HAPPEN
				// if this fires revaluate how im marshaling and unmarshling json
				// -----------------
				// msg := fmt.Sprintf("could not receive processed file in transit: %s", err.Error())
				// extractionResult.FailedFiles[shared.FileID{}] = msg
				panic(err)
			} // if
			fileToAdd := shared.FileID{
				Title:     string(success.Title),
				ObjectKey: string(success.ObjectKey),
				URL:       string(success.URL),
			}
			// add successfully converted file
			// extractionResult.SuccessFiles[fileToAdd] = []byte(success.Body)
			l.log(fileToAdd, success.Body)
		} // for
	}

	go readStdout()

	readStderr := func() {
		wg.Done()
		// Read errors from std err
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			var failure FileJSON
			_ = json.Unmarshal(scanner.Bytes(), &failure)
			fileToAdd := shared.FileID{
				ObjectKey: string(failure.Err),
			}
			// extractionResult.FailedFiles[shared.FileID{}] = string(msg)
			l.log(fileToAdd, string(failure.Err))
		} // for
	} // Routine

	go readStderr()

	// write json objects
	for id, r := range files {

		data, err := buildJSON(&l, id, r)
		if err != nil {
			continue
		} // if
		stdin.Write(data)
		stdin.Write([]byte("\n"))
	} // for
	stdin.Close()

	wg.Wait()
	cmd.Wait()
	return &l.e, err
} // Extract

/*
Returns a properly encoded json object ready for streaming

On error, populates ExtractionResult.FailedFiles
*/
func buildJSON(l *Locker, id shared.FileID, r io.ReadCloser) ([]byte, error) {
	var err error
	var buf []byte
	buf, err = io.ReadAll(r)
	if err != nil {
		msg := fmt.Sprintf("could not read pdf %s: %s", id.ObjectKey, err.Error())
		// e.FailedFiles[id] = msg
		l.log(id, msg)
		return buf, err
	} // if

	// encode body in b64
	encodedBuf := base64.StdEncoding.EncodeToString(buf)

	// encode data
	raw := FileJSON{
		Body:      encodedBuf,
		Title:     id.Title,
		ObjectKey: id.ObjectKey,
		URL:       id.URL,
	}
	data, err := json.Marshal(raw)
	if err != nil {
		msg := fmt.Sprintf("could not encode json payload %s: %s", id.ObjectKey, err.Error())
		// e.FailedFiles[id] = msg
		l.log(id, msg)
		return data, err
	} // if
	return data, err
} // buildJSON
