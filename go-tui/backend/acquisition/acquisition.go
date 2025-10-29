/*
Handles the orchestration of the acquisition flow
*/
package acquisition

import (
	"errors"
	"fmt"
	"tui/backend/types"
	"tui/backend/types/shared"
)

// represents an acquisition instance
type Acquisition struct {
	shared.Workflow
}

// Executes the acquisition flow
func (a *Acquisition) Run(input types.AcquisitionInput) (types.AcquisitionResult, error) {
	var err error
	result := types.AcquisitionResult{
		SuccessFiles: make([]shared.File, 0),
		FailedFiles:  make([]shared.FailedFile, 0),
		Log:          make([]string, 0),
	}

	// initialize clients
	if err := a.InitializeClients(); err != nil {
		return result, err
	} // if
	defer a.PGClient.Close()

	// 1) scrape urls
	scrapeResult, err := Scrape(input.Query, input.FileCount)
	if err != nil {
		err := errors.New("could not visit url")
		return result, err
	}
	// log count of pages searched to result and result log
	result.PagesSearched = scrapeResult.PageCount
	pagesSearchedMsg := fmt.Sprintf("%d out of a maximum %d pages worth of results scraped", result.PagesSearched, MaxPages)
	result.Log = append(result.Log, pagesSearchedMsg)

	// log count of file urls scraped to result log
	result.URLsScraped = scrapeResult.URLCount
	urlCountMsg := fmt.Sprintf("%d out of a requested %d URLs scraped", scrapeResult.URLCount, input.FileCount)
	result.Log = append(result.Log, urlCountMsg)

	// 2) download files
	urlsToDownload := scrapeResult.URLMap
	downloadResults, err := a.DownloadFiles(urlsToDownload)
	if err != nil {
		return result, err
	} // if
	// log results of downloaded files
	for fileInfo, _ := range downloadResults.SuccessFiles {
		title, url := fileInfo.Title, fileInfo.URL
		successMsg := fmt.Sprintf("File [%s] downloaded successfully from [%s]", title, url)
		result.Log = append(result.Log, successMsg)
	} // for

	// log total amount of downloaded files
	downloadedFileCount := len(downloadResults.SuccessFiles)
	successMsg := fmt.Sprintf("%d files downloaded successfully out of %d URLs", downloadedFileCount, result.URLsScraped)
	result.Log = append(result.Log, successMsg)

	// log results of failed files
	for fileInfo, report := range downloadResults.FailedFiles {
		title, url := fileInfo.Title, fileInfo.URL
		failedMsg := fmt.Sprintf("%s: Could not download %s from %s", report, title, url)
		result.Log = append(result.Log, failedMsg)
	} // for

	// log total amount of failed files
	failedFileCount := len(downloadResults.FailedFiles)
	failedMsg := fmt.Sprintf("Failed to download %d files", failedFileCount)
	result.Log = append(result.Log, failedMsg)

	// log total amount of skipped files
	skippedMsg := fmt.Sprintf("%d files skipped - already exist in database", downloadResults.ExistingFilesCount)
	result.Log = append(result.Log, skippedMsg)

	// 3) upload to minio
	uploadResult, err := a.Upload(&downloadResults.SuccessFiles)
	if err != nil {
		return result, err
	} // if

	// 5) populate AcquisitionResult.Log

	// total count of successfully downloaded and logged files
	successfulCount := uploadResult.SuccessUploadCount
	// total count of files requested by user
	requestedFileCount := input.FileCount

	// add final log messages
	result.Log = append(result.Log, fmt.Sprintf("%d file(s) requested", requestedFileCount))
	existingCount := downloadResults.ExistingFilesCount + uploadResult.ExistingFilesCount
	existingMsg := fmt.Sprintf("%d files already exist in database", existingCount)
	result.Log = append(result.Log, existingMsg)
	result.Log = append(result.Log, fmt.Sprintf("%d file(s) downloaded, stored, and logged in database", successfulCount))

	// add files to acquisition object
	sfiles := uploadResult.SuccessFiles
	ffiles := uploadResult.FailedFiles
	result.SuccessFiles = append(result.SuccessFiles, sfiles...)
	result.FailedFiles = append(result.FailedFiles, ffiles...)

	// 6) return AcquisitionResult
	return result, err
} // Run
