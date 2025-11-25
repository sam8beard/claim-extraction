package main

import (
	"context"
	"fmt"
	"github.com/sam8beard/claim-extraction/internal/types"
	"github.com/sam8beard/claim-extraction/internal/workflows/acquisition"
	"github.com/sam8beard/claim-extraction/internal/workflows/conversion"
	"github.com/sam8beard/claim-extraction/internal/workflows/processing"
	"log"
	"os"
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
	fc, _ := strconv.Atoi(os.Args[2])
	fmt.Printf("\n\t\tQuery: %s | File Count: %d\n\n", q, fc)

	acqInput := types.AcquisitionInput{
		Query:     q,
		FileCount: fc,
	}
	log.Println("\n\t\t------------ Starting acquisition... ------------ \n\n")
	acqResult, err := a.Run(ctx, acqInput)
	if err != nil {
		log.Fatalf("acquisition error: %v\n", err)
	} // if
	conInput := types.ConversionInput(acqResult)

	log.Println("\n\t\t------------ Starting conversion... ------------ \n\n")
	conResult, err := c.Run(ctx, &conInput)
	if err != nil {
		log.Fatalf("conversion error: %v\n", err)
	} // if
	procInput := types.ProcessingInput{
		ConvertedFiles: conResult.ConvertedFiles,
	}
	log.Println("\n\t\t------------ Starting processing... ------------ \n\n")
	time.Sleep(time.Second * 3)
	procResult, err := p.Run(ctx, &procInput)
	if err != nil {
		log.Fatalf("processing error: %v\n", err)
	} // if

	// printNLP(procResult)

	resultHeader := "[SUCCESS]"
	resultMsg := "pipeline executed successfully"
	if len(procResult.FileData) == 0 {
		resultHeader = "[FAIL]"
		resultMsg = "failed to process any files"
	} // if
	fmt.Printf("\n\n------------------------- %s -------------------------\n", resultHeader)
	fmt.Printf("\t\t%s\n\n", resultMsg)
	fmt.Println("-------------------------- [REPORT] ---------------------------")
	fmt.Printf("\t\t      %d files requested\n", fc)

	fmt.Printf("\t\t      %d files retrieved\n", len(acqResult.SuccessFiles))

	fmt.Printf("\t\t      %d files converted to text\n", len(conResult.ConvertedFiles))

	fmt.Printf("\t\t      %d files processed\n", len(procResult.FileData))

} // main

func printNLP(nr *processing.NLPResult) {
	fileData := nr.FileData
	for _, data := range fileData {
		fmt.Printf("\nSpan data for %s ---------------------------\n", data.FileName)
		for _, spanData := range data.ClaimSpans {
			fmt.Printf("\nType: %s\n", spanData.Type)
			fmt.Printf("\nText: %s\n", spanData.Text)
			fmt.Printf("\nSent: %s\n", spanData.Sent)
			fmt.Printf("\nConfidence: %f\n", spanData.Confidence)
		} // for
		fmt.Println("\nFile data ---------------------------")
		fmt.Printf("\nObject Key: %s\n", data.ObjectKey)
		fmt.Printf("\nFile Name: %s\n", data.FileName)
		fmt.Printf("\nClaim Score: %f\n", data.ClaimScore)
	} // for
} // printNLPResult
