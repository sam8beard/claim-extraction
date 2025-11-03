package conversion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

// represents a conversion instance
type Conversion struct {
	shared.Workflow
}

func (c *Conversion) Run(ctx context.Context, input types.ConversionInput) (types.ConversionResult, error) {
	var err error
	finalResult := types.ConversionResult{
		ConvertedFiles: make([]shared.File, 0),
		Log:            make([]string, 0),
	}

	// initialize clients
	if err := c.InitializeClients(); err != nil {
		return finalResult, err
	} // if
	defer c.PGClient.Close()

	/* Outline conversion process here */

	log.Println("Beginning download...")
	dResult, err := c.Download(ctx, input)
	if err != nil {
		log.Println("Error on download")

		// this shouldnt be a problem
		err = errors.New("could not download files")
		return finalResult, err
	} // if
	log.Println("Download successful...")

	log.Println("Beginning extraction...")
	eResult, err := c.Extract(ctx, dResult)
	if err != nil {
		log.Println("Error on extraction")

		err = errors.New("could not extract files")
		return finalResult, err
	} // if
	log.Println("Extraction successful")

	eFiles := eResult.SuccessFiles

	log.Println("Beginning upload...")
	upResult, err := c.Upload(ctx, &eFiles)
	if err != nil {
		log.Println("Error on upload")

		err = errors.New("could not upload files")
		return finalResult, err
	} // if
	log.Println("Upload successful")

	var failedLog string
	var successLog string

	for _, fFile := range upResult.FailedFiles {
		failedLog += fmt.Sprintf("\n%s\n", fFile.Report)
	} // for
	finalResult.Log = append(finalResult.Log, failedLog)

	for _, sFile := range upResult.SuccessFiles {
		successLog += fmt.Sprintf("\n%s: %s\n", sFile.ObjectKey, sFile.Status)
		finalResult.ConvertedFiles = append(finalResult.ConvertedFiles, sFile)
	} // for
	finalResult.Log = append(finalResult.Log, successLog)

	// What other files and functions do we need?
	// use the keys from ConversionInput to download files from Minio
	// downloadResult, err := c.Download(ctx, input)
	// if err != nil {
	// 	return finalResult, err
	// } // if

	// pass file readers to function that calls python sub process
	// (apparently no reliable library for pdf conversion in go)
	// extractionResult, err := c.Extract(ctx, downloadResult)
	// if err != nil {
	// 	return finalResult, err
	// } // if

	// // once text is parsed, update text_extracted in documents table
	// updateResult, err := c.Update(ctx, extractionResult)
	// if err != nil {
	// 	return finalResult, err
	// } // if

	// // make new object key, and upload parsed text to bucket
	// uploadResult, err := c.Upload(ctx, extractionResult)
	// if err != nil {
	// 	return finalfinalResult, err
	// } // if

	// log all results

	// return final conversion result
	return finalResult, err
} // Run
