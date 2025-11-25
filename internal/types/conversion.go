package types

import "github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"

/*
Acquisition types
*/

// this is the same as AcquisitionOutput
type ConversionInput struct {
	SuccessFiles  []shared.File       // list of successfully aquired files
	FailedFiles   []shared.FailedFile // list of failed files
	Log           []string            // report of acquisition
	PagesSearched int                 // amount of pages searched
	URLsScraped   int                 // amount of urls scraped
}

type ConversionResult struct {
	ConvertedFiles []shared.File
	Log            []string
}
