package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Result struct {
	Name string `json:"name,omitempty"`
	Job  string `json:"job,omitempty"`
}

func PipePython() ([]Result, error) {
	var err error
	results := make([]Result, 0)

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

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	// start python program
	if err := cmd.Start(); err != nil {
		panic(err)
	} // if
	// write to stdin
	stuff := []Result{
		{
			Name: "Sammy",
			Job:  "Welder",
		},
		{
			Name: "Jason",
		},
	}

	for _, item := range stuff {
		fmt.Printf("%v", item)
		json_item, _ := json.Marshal(item)
		if _, err := stdin.Write(json_item); err != nil {
			panic(err)
		} // if
		_, err = stdin.Write([]byte("\n"))

	} // for

	stdin.Close()

	// read from stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var output Result
		if err := json.Unmarshal(scanner.Bytes(), &output); err != nil {
			fmt.Print(err)
		} // if
		results = append(results, output)
	} // for
	cmd.Wait()
	return results, err
} // TestStdin
