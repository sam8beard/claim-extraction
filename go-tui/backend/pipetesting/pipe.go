package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Result struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type Response struct {
	SuccessFiles map[string]string
	FailedFiles  map[string]string
}

type Locker struct {
	mu sync.Mutex
	r  Response
}

func (l *Locker) log(name string, id string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch id {
	case "1234":
		l.r.SuccessFiles[name] = id
	case "5432":
		l.r.FailedFiles[name] = id

	}
}

func PipePython() (Response, error) {
	var err error
	results := make([]Result, 0)

	l := Locker{
		r: Response{
			SuccessFiles: make(map[string]string),
			FailedFiles:  make(map[string]string),
		},
	}
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
	scriptPath := filepath.Join(pythonDir, "testing.py")
	// pythonVenv := "python/venv/bin/python3.12"
	cmd := exec.Command(pythonExec, "-u", scriptPath)
	cmd.Dir = pythonDir

	// copy curr environment, but inject venv info
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", venvDir),
		fmt.Sprintf("PATH=%s%c%s", filepath.Join(venvDir, "bin"), os.PathListSeparator, os.Getenv("PATH")),
	)

	stuff := []Result{
		{
			Name: "Sammy",
			ID:   "1234",
		},
		{
			Name: "Jason",
			ID:   "5432",
		},
	}

	var wg sync.WaitGroup
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	// start python program
	if err := cmd.Start(); err != nil {
		panic(err)
	} // if

	wg.Add(1)
	readStdout := func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var output Result
			if err := json.Unmarshal(scanner.Bytes(), &output); err != nil {
				fmt.Print(err)
			} // if
			l.log(output.Name, output.ID)
			results = append(results, output)

		} // for
	}

	go readStdout()

	for _, item := range stuff {
		// fmt.Printf("%v", item)
		json_item, _ := json.Marshal(item)
		if _, err := stdin.Write(json_item); err != nil {
			panic(err)
		} // if
		_, err = stdin.Write([]byte("\n"))

	} // for

	stdin.Close()

	wg.Wait()
	cmd.Wait()
	return l.r, err
} // TestStdin
