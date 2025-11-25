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
}

func (p *Processing) NLP(ctx context.Context, f *FetchResult) (*NLPResult, error) {
	result := NLPResult{
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
	_, currentFile, _, _ := runtime.Caller(0) // file where this code is
	currentDir := filepath.Dir(currentFile)
	pythonDir := filepath.Join(currentDir, "python")
	venvDir := filepath.Join(pythonDir, "venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "nlp_processing.py")

	cmd := exec.Command(pythonExec, "-u", "-W", "ignore", scriptPath)
	cmd.Dir = pythonDir

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
	log.Println("Executing script...")
	if err := cmd.Start(); err != nil {
		msg := fmt.Errorf("start script error: %v", err)
		return nil, msg
	}

	// add routines to waitgroup
	var wg sync.WaitGroup
	wg.Add(2)

	readStdout := func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var fd FileData
			if err := json.Unmarshal(scanner.Bytes(), &fd); err != nil {
				log.Printf("bad JSON from python: %v", err)
				continue
			} // if
			result.FileData = append(result.FileData, fd)
			log.Printf("read file successfully: %s", fd.ObjectKey)
		} // for
	} // readStdout

	log.Println("Begin reading from stdout...")
	// start routine to read from stdout
	go readStdout()

	readStderr := func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			var fd FileData
			line := scanner.Bytes()
			if err := json.Unmarshal(line, &fd); err != nil {
				log.Printf("bad JSON from stderr: %v", err)
				continue
			} // if
			log.Printf("proccessing error for file: %s: %s", fd.ObjectKey, fd.Error)
		} // for
		if err := scanner.Err(); err != nil && err != io.EOF {
			log.Printf("stderr scanner error: %v", err)
		}
	} // readStderr

	log.Println("Began reading from stderr...")
	go readStderr()

	// defer closing of stdin and ignore error
	for file, reader := range files {
		// get file name from object key for json data
		fileName := file.ObjectKey
		fileName = strings.Replace(fileName, "processed/", "", 1)
		jsonObject := PythonInput{
			FileName:  fileName,
			ObjectKey: file.ObjectKey,
			Error:     "",
			//	Content:   content.String(), // we need to stream the body instead
		}

		// marshal metadata
		metaJson, err := json.Marshal(jsonObject)
		if err != nil {
			log.Printf("marshal error: %v", err)
		} // if

		// write metadata + newline
		log.Printf("writing metadata for %s...\n", fileName)
		if _, err := stdin.Write(metaJson); err != nil {
			msg := fmt.Sprintf("write meta to stdin: %v", err)
			log.Println(msg)
			continue
		} // if
		if _, err := stdin.Write([]byte("\n")); err != nil {
			msg := fmt.Sprintf("write newline to stdin: %v", err)
			log.Println(msg)
			continue
		} // if

		// wrap stdin with b64 encoder
		encoder := base64.NewEncoder(base64.StdEncoding, stdin)

		log.Printf("writing file %s", fileName)
		// stream file to python
		_, err = io.Copy(encoder, reader)

		// close encoder
		if cerr := encoder.Close(); cerr != nil {
			log.Printf("encoder close error: %v", cerr)
		} // if

		// check error on copy
		if err != nil {
			msg := fmt.Sprintf("IO copy to python subprocess: %v", err)
			log.Println(msg)
			continue
		} // if

		// write file delimeter
		if _, err := stdin.Write([]byte(bodySentinel)); err != nil {
			msg := fmt.Errorf("write sentinel: %v", err)
			log.Println(msg)
			continue
		} // if

		// close reader
		if err := reader.Close(); err != nil {
			msg := fmt.Sprintf("close reader: %v", err)
			log.Println(msg)
			return nil, err
		} // if

	} // for

	log.Println("Begin writing to stdout...")

	if err := stdin.Close(); err != nil {
		msg := fmt.Errorf("stdin close error: %v", err)
		return nil, msg
	} // if

	// wait for routines to finish and script to finish executing
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		msg := fmt.Errorf("cmd.Wait error: %v", err)
		return nil, msg
	} // if
	return &result, nil
} // NLP
