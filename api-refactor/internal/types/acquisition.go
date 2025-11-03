package types

import "github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"

/*
Acquisition types
*/

type AcquisitionInput struct {
	Query     string
	FileCount int
}

type AcquisitionResult struct {
	SuccessFiles  []shared.File       // list of successfully aquired files
	FailedFiles   []shared.FailedFile // list of failed files
	Log           []string            // report of acquisition
	PagesSearched int                 // amount of pages searched
	URLsScraped   int                 // amount of urls scraped
}
