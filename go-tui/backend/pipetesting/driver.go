package main

import (
	"fmt"
)

func main() {
	output, err := PipePython()
	if err != nil {
		fmt.Println(err)
	} // if

	succFiles := output.SuccessFiles
	failFiles := output.FailedFiles

	fmt.Printf("\nSuccess Files: %v\n", succFiles)

	fmt.Printf("\nFailed Files: %v\n", failFiles)
}
