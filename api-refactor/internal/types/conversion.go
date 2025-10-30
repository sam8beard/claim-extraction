package types

import (
	"tui/backend/types/shared"
)

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
