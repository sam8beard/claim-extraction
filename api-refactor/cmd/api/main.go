package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/acquisition"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/conversion"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/workflows/processing"
)

func main() {

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
		Query:     "what does ai have to do with politics",
		FileCount: 2,
	}
	fmt.Println("Starting acquisition")
	acqResult, err := a.Run(ctx, acqInput)
	if err != nil {
		panic(err)
	} // if

	conInput := types.ConversionInput{
		SuccessFiles:  acqResult.SuccessFiles,
		FailedFiles:   acqResult.FailedFiles,
		Log:           acqResult.Log,
		PagesSearched: acqResult.PagesSearched,
		URLsScraped:   acqResult.URLsScraped,
	}
	fmt.Println("Starting conversion")

	conResult, err := c.Run(ctx, conInput)
	if err != nil {
		panic(err)
	} // if
	log.Println("firing")

	procInput := types.ProcessingInput{
		ConvertedFiles: conResult.ConvertedFiles,
	}
	fmt.Println("Starting processing")

	procResult, err := p.Run(ctx, &procInput)

	processing.PrintNLPResult(procResult)
	fmt.Printf("%v", procResult.FileData)

} // main
