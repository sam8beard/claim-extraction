package main

import (
	"fmt"
)

func main() {
	output, err := PipePython()
	if err != nil {
		fmt.Println(err)
	} // if

	for _, result := range output {
		fmt.Printf("%s\n", result)
	} // for
}
