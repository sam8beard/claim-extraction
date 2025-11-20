package main

import (
	"context"
	"fmt"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/acquisition"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/conversion"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/processing"
	"log"
	"runtime/debug"
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
	acqInput := types.AcquisitionInput{
		Query:     "short quick summary on ai and climate",
		FileCount: 3,
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
	printNLP(*procResult)
	fmt.Printf("%v", procResult.FileData)

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
