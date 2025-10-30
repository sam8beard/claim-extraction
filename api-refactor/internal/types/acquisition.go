package types

import "tui/backend/types/shared"

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
