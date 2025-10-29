package conversion

import (
	"tui/backend/types"
	"tui/backend/types/shared"
)

// represents an acquistion instance
type Conversion struct {
	shared.Workflow
}

func (c *Conversion) Run(input types.ConversionInput) (types.ConversionResult, error) {
	var err error
	result := types.ConversionResult{
		ConvertedFiles: make([]shared.File, 0),
		Log:            make([]string, 0),
	}

	// initialize clients
	if err := c.InitializeClients(); err != nil {
		return result, err
	} // if
	defer c.PGClient.Close()

	// Outline conversion process here
	/*

		// What other files and functions do we need?


	*/
	//
	return result, err
} // Run
