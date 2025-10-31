package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func PipePython() ([][]byte, error) {
	var err error
	results := make([][]byte, 0)
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
	stuff := []string{
		"thing 1",
		"thing 2",
		"thing 3",
		"thing 4",
	}

	for _, item := range stuff {
		if _, err := stdin.Write([]byte(item)); err != nil {
			panic(err)
		} // if
	} // for

	stdin.Close()

	// read from stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		results = append(results, scanner.Bytes())
	} // for
	cmd.Wait()

	return results, err
} // TestStdin
