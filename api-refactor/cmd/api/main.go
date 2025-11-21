package main

import (
	"context"
	"fmt"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/acquisition"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/conversion"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/processing"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

func main() {
	// configure logger
	log.SetFlags(log.Lshortfile)
	// initialize workflows
	a := acquisition.Acquisition{}
	c := conversion.Conversion{}
	p := processing.Processing{}

	if err := a.InitializeClients(); err != nil {
		panic(err)
	} // if
	if err := c.InitializeClients(); err != nil {
		panic(err)
	} // if
	if err := p.InitializeClients(); err != nil {
		panic(err)
	} // if

	ctx := context.Background()

	// cli args for testing
	q := string(os.Args[1])
	fc, err := strconv.Atoi(os.Args[2])
	log.Printf("Query: %s --- File Count: %d\n", q, fc)

	acqInput := types.AcquisitionInput{
		Query:     q,
		FileCount: fc,
	}
	fmt.Println("Starting acquisition")
	acqResult, err := a.Run(ctx, acqInput)
	if err != nil {
		panic(err)
	} // if
	// log.Fatalf("Acq result: %s, %d, %d, %d, %v", acqInput.Query, acquisition.MaxPages, acqResult.URLsScraped, acqResult.PagesSearched, acqResult.SuccessFiles)
	conInput := types.ConversionInput{
		SuccessFiles:  acqResult.SuccessFiles,
		FailedFiles:   acqResult.FailedFiles,
		Log:           acqResult.Log,
		PagesSearched: acqResult.PagesSearched,
		URLsScraped:   acqResult.URLsScraped,
	}
	fmt.Println("Starting conversion")

	conResult, err := c.Run(ctx, conInput) // NOTE: never returning
	if err != nil {
		log.Fatal(string(debug.Stack()))
	} // if
	log.Println("firing")
	//	log.Fatalf("Con result: %v", conResult.Log)
	procInput := types.ProcessingInput{
		ConvertedFiles: conResult.ConvertedFiles,
	}
	fmt.Println("Starting processing")
	time.Sleep(time.Second * 3)
	procResult, err := p.Run(ctx, &procInput)
	if err != nil {
		log.Println("firing in err block for p.Run")
		log.Println(string(debug.Stack()))
		log.Printf("error on p.Run: %v\n", err)
	} // if
	//printNLP(*procResult)
	//fmt.Printf("%v", procResult.FileData)
	fmt.Println("------------------------- [SUCCESS] -------------------------")
	fmt.Println("\t\tpipeline executed successfully\n\n")
	fmt.Println("-------------------------- [INFO] ---------------------------")
	fmt.Printf("\t\t\t%d files processed\n", len(procResult.FileData))
	fmt.Printf("\t\t\t%d files requested\n", fc)
} // main

func printNLP(nr processing.NLPResult) {
	fileData := nr.FileData
	for _, data := range fileData {
		fmt.Println("\nFile data ---------------------------")
		fmt.Printf("\nObject Key: %s\n", data.ObjectKey)
		fmt.Printf("\nFile Name: %s\n", data.FileName)
		fmt.Printf("\nClaim Score: %f\n", data.ClaimScore)
		fmt.Printf("\nSpan data for %s ---------------------------\n", data.FileName)
		for _, spanData := range data.ClaimSpans {
			fmt.Printf("\nType: %s\n", spanData.Type)
			fmt.Printf("\nText: %s\n", spanData.Text)
			fmt.Printf("\nSent: %s\n", spanData.Sent)
			fmt.Printf("\nConfidence: %f\n", spanData.Confidence)
		} // for
	} // for
} // printNLPResult
