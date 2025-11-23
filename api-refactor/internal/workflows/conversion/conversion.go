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

func (c *Conversion) Run(ctx context.Context, input *types.ConversionInput) (*types.ConversionResult, error) {
	var err error
	finalResult := types.ConversionResult{
		ConvertedFiles: make([]shared.File, 0),
		Log:            make([]string, 0),
	}

	log.Println("Beginning download...")
	dResult, err := c.Download(ctx, input)
	log.Println("Firing right after download")
	if err != nil {
		log.Println("Error on download")
		// this shouldnt be a problem
		err = errors.New("could not download files")
		return nil, err
	} // if
	log.Println("Download successful...")
	//log.Printf("Result of download: %v\n", dResult.SuccessFiles)
	for key, value := range dResult.SuccessFiles {
		if value == nil {
			log.Printf("READER DOES NOT EXIST FOR %s\n", key.OriginalKey)
			continue
		} // if
		//log.Printf("\nOriginal Key: %s\nValue: %v\n\n", key.OriginalKey, value)
	} // for
	log.Println("Beginning extraction...")
	//	log.Fatalf("num of downloaded files: %d", len(dResult.SuccessFiles))
	eResult, err := c.Extract(ctx, dResult)
	if err != nil {
		log.Println("Error on extraction")
		log.Fatal("firing on extraction error")
		err = errors.New("could not extract files")
		return nil, err
	} // if

	// no files were converted
	if len(eResult.SuccessFiles) == 0 {
		if len(eResult.FailedFiles) != 0 {
			log.Println("Only failed files returned from Extract")
			for file, msg := range eResult.FailedFiles {
				log.Printf("%s -- %s\n", file.ObjectKey, msg)
			} // for
		} // if
		log.Fatal("No successful files were returned from Extract")
	} // if

	log.Println("Extraction successful")

	eFiles := eResult.SuccessFiles

	log.Println("Beginning upload...")

	upResult, err := c.Upload(ctx, eFiles)
	for _, file := range upResult.SuccessFiles {
		log.Printf("%s %s", file.FileName, file.ObjectKey)
	} // for
	if err != nil {
		log.Println("Error on upload")

		err = errors.New("could not upload files")
		return nil, err
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

	// return final conversion result
	return &finalResult, err
} // Run
