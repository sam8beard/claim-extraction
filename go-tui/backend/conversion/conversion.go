package conversion

import (
	"tui/backend/types"
	"tui/backend/types/shared"
)

// represents a conversion instance
type Conversion struct {
	shared.Workflow
}

func (c *Conversion) Run(input types.ConversionInput) (types.ConversionResult, error) {
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
	pdfResult, err := c.GetPDFs()
	// use the keys from ConversionInput to downlaod files from Minio

	// pass file readers to python sub process
	// (apparently no reliable library for pdf conversion in go)

	// once text is parsed, update text_extracted in documents table

	// make new object key, and upload parsed text to bucket

	//

	//
	return finalResult, err
} // Run
