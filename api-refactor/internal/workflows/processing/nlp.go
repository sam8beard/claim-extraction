package processing

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type PythonInput struct {
	ObjectKey string `json:"objectKey"`
	Content   string `json:"content"`
	FileName  string `json:"fileName"`
}

type FileData struct {
	ObjectKey  string      `json:"objectKey"`
	FileName   string      `json:"fileName"`
	ClaimScore float64     `json:"claimScore"`
	ClaimSpans []ClaimSpan `json:"claimSpans"`
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
	var err error
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
	// projectRoot, _ := os.Getwd()
	// pythonDir := filepath.Join(projectRoot, "python")
	// venvDir := filepath.Join(pythonDir, "venv")
	// pythonExec := filepath.Join(venvDir, "bin", "python3")
	// scriptPath := filepath.Join(pythonDir, "nlp_processing.py")
	// cmd := exec.Command(pythonExec, "-u", scriptPath)
	// cmd.Dir = pythonDir

	// // copy curr environment, but inject venv info
	// cmd.Env = append(os.Environ(),
	// 	fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
	// 	fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	// )
	_, currentFile, _, _ := runtime.Caller(0) // file where this code is
	currentDir := filepath.Dir(currentFile)
	pythonDir := filepath.Join(currentDir, "python")
	venvDir := filepath.Join(pythonDir, "venv")
	pythonExec := filepath.Join(venvDir, "bin", "python3")
	scriptPath := filepath.Join(pythonDir, "convert_pdf.py")

	cmd := exec.Command(pythonExec, "-u", scriptPath)
	cmd.Dir = pythonDir

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
		fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	)

	// open pipes
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	fmt.Println("Executing script...")
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// add routines to waitgroup
	var wg sync.WaitGroup
	wg.Add(2)

	readStdout := func() {
		defer wg.Done()
		defer stdout.Close()
		fmt.Println("firing inside read")
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var fd FileData
			if err := json.Unmarshal(scanner.Bytes(), &fd); err != nil {
				log.Printf("bad JSON from python: %v", err)
				continue
			} // if
			result.FileData = append(result.FileData, fd)
		} // for

	} // readStdout

	fmt.Println("Begin reading from stdout...")
	// start routine to read from stdout
	go readStdout()

	writeStdin := func() {
		fmt.Println("firing inside write")
		defer wg.Done()
		defer stdin.Close()
		fmt.Println("firing inside write")
		for file, content := range files {
			// get file name from object key for json data
			fileName := file.ObjectKey
			fileName = strings.Replace(fileName, "processed/", "", 1)
			input := PythonInput{
				FileName:  fileName,
				ObjectKey: file.ObjectKey,
				Content:   content.String(),
			}
			data, err := json.Marshal(input)
			if err != nil {
				log.Printf("marshal error: %v", err)
			} // if
			fmt.Printf("writing file %s...\n", fileName)
			stdin.Write(data)
			stdin.Write([]byte("\n"))
		} // for
	} // writeStdin

	fmt.Println("Begin writing to stdout...")
	// start routine to write to stdin
	go writeStdin()
	// wait for routines to finish and script to finish executing
	wg.Wait()
	cmd.Wait()
	return &result, err
} // NLP
