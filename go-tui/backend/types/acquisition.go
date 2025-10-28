package types

/*
Acquisition types
*/

type AcquisitionInput struct {
	Query     string
	FileCount int
}

type File struct {
	FileName string // name of file
	Key      string // object key of file
	URL      string // url of file
	Status   string // status of file (downloaded, uploaded, failed)
}

type FailedFile struct {
	URL    string // url of file
	Report string // reason for why file failed in acquisition
}

type AcquisitionResult struct {
	SuccessFiles  []File       // list of successfully aquired files
	FailedFiles   []FailedFile // list of failed files
	Log           []string     // report of acquisition
	PagesSearched int          // amount of pages searched
	URLsScraped   int          // amount of urls scraped
}

// type Scraper struct {

// }
