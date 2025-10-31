package conversion

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (l *Locker) log(f shared.FileID, u any) {
	l.mu.Lock() // calling routine blocks other routines from modifying the mutex
	defer l.mu.Unlock()
	log.Fatalf("Firing %v\t and %v", f, u)

	switch val := u.(type) {
	case string:
		// this is a failed file
		l.e.FailedFiles[f] = val
	case []byte:
		// this is a successful file
		l.e.SuccessFiles[f] = val
	} // switch
}

// func reader(scanner *bufio.Scanner) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	for scanner.Scan() {
// 		fmt.Println()
// 	}

// }
func (c *Conversion) Extract(ctx context.Context, d shared.DownloadResult) (ExtractionResult, error) {
	var err error
	l := Locker{
		e: ExtractionResult{
			FailedFiles:  make(map[shared.FileID]string),
			SuccessFiles: make(map[shared.FileID][]byte),
		},
	}

	var wg sync.WaitGroup

	files := d.SuccessFiles
	cmd := exec.Command("python3", "u", "python/convert_pdf.py")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// // initialize scanners for each pipe
	// scanout := bufio.NewScanner(stdout)
	// scanerr := bufio.NewScanner(stderr)

	if err := cmd.Start(); err != nil {
		panic(err)
	} // if

	// tell the waitgroup we're waiting for two routines to finish
	// wg.Add(2)
	readStdout := func() {
		wg.Add(1)
		defer wg.Done()

		// Reading converted files line by line
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var success FileJSON
			if err := json.Unmarshal(scanner.Bytes(), &success); err != nil {
				log.Printf("error decoding response: %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("syntax error at byte offset %d, %v", e.Offset, success)
				} // if
				log.Printf("success block: %v", success)

			}

			fileToAdd := shared.FileID{
				Title:     success.Title,
				ObjectKey: success.ObjectKey,
				URL:       success.URL,
			}
			log.Printf("%v", fileToAdd)
			// add successfully converted file
			l.log(fileToAdd, success.Body)
		} // for
	}

	go readStdout()

	readStderr := func() {
		// Read errors from std err
		wg.Add(1)
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			var failure FileJSON
			if err := json.Unmarshal(scanner.Bytes(), &failure); err != nil {
				log.Printf("error decoding response: %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("syntax error at byte offset %d, %v", e.Offset, failure)
				} // if
				log.Printf("failure block: %v", failure)
			}

			fileToAdd := shared.FileID{
				ObjectKey: string(failure.Err),
			}

			// extractionResult.FailedFiles[shared.FileID{}] = string(msg)
			l.log(fileToAdd, string(failure.Err))
		} // for
	} // Routine

	go readStderr()
	wg.Wait()

	// write json objects
	for id, r := range files {
		raw, err := buildJSON(&l, id, r)
		// log.Fatalf("%v", &data)
		if err != nil {
			continue
		} // if

		data, err := json.Marshal(raw)
		if err != nil {
			panic(err)
		}
		stdin.Write(data)
		stdin.Write([]byte("\n"))
		// enc := json.NewEncoder(stdin)
		// if err := enc.Encode(&data); err != nil {
		// 	panic(err)
		// }
		// if err := enc.Encode("\n"); err != nil {
		// 	panic(err)
		// }
	} // for
	stdin.Close()

	// wg.Wait()
	cmd.Wait()

	return l.e, err
} // Extract

/*
Returns a properly encoded json object ready for streaming

On error, populates ExtractionResult.FailedFiles
*/
func buildJSON(l *Locker, id shared.FileID, r io.ReadCloser) (FileJSON, error) {
	var err error
	var buf []byte
	var raw FileJSON
	buf, err = io.ReadAll(r) // ERROR HERE
	if err != nil {
		msg := fmt.Sprintf("could not read pdf %s: %s", id.ObjectKey, err.Error())
		// e.FailedFiles[id] = msg
		l.log(id, msg)
		return FileJSON{}, err
	} // if

	// encode body in b64
	encodedBuf := base64.StdEncoding.EncodeToString(buf)

	// encode data
	raw = FileJSON{
		Body:      encodedBuf,
		Title:     id.Title,
		ObjectKey: id.ObjectKey,
		URL:       id.URL,
		Err:       "",
	}
	// data, err := json.Marshal(raw)
	// if err != nil {
	// 	msg := fmt.Sprintf("could not encode json payload %s: %s", id.ObjectKey, err.Error())
	// 	// e.FailedFiles[id] = msg
	// 	l.log(id, msg)
	// 	return data, err
	// } // if
	return raw, err
} // buildJSON
