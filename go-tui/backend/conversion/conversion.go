package conversion

import (
	"context"
	"tui/backend/types"
	"tui/backend/types/shared"
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
