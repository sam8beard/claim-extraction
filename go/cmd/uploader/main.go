package main

import ( 
	"fmt"
	"flag"
)

func main() { 
	var filePath string
	var source string 

	// identify which flags to look for - user must enter file path and source
	flag.StringVar(&filePath, "file", "", "Path to the document")
	flag.StringVar(&source, "source", "", "Source of the document")

	// parse flags provided 
	flag.Parse()

	if filePath == "" || source == "" { 
		fmt.Println("Must provide file path and source to use")
		return
	} // if 

	fmt.Println(filePath, source)

} // main