package acquisition

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tui/backend/types/shared"
)

type DownloadResults struct {
	FailedFiles        map[shared.FileID]string
	SuccessFiles       map[shared.FileID]io.ReadCloser
	DownloadCount      int
	ExistingFilesCount int
}

// THIS IS WHERE WE DETERMINE WHAT FILES ARE ALREADY IN OUR DATABASE
// WHAT SHOULD WE DO?
// checks if file already exists in database
func (a *Acquisition) checkFile(ctx context.Context, file shared.FileID) bool {
	var exists bool
	title := fmt.Sprint(file.Title, ".pdf")
	query := `
	SELECT EXISTS (SELECT 1 FROM documents WHERE file_name = $1)
	`
	_ = a.PGClient.QueryRow(ctx, query, title).Scan(&exists)

	return exists
} // checkFile

func (a *Acquisition) Download(ctx context.Context, urlMap map[string]string) (DownloadResults, error) {
	var err error
	results := DownloadResults{
		FailedFiles:  make(map[shared.FileID]string),
		SuccessFiles: make(map[shared.FileID]io.ReadCloser, 0),
	}
	for title, url := range urlMap {
		fileKey := shared.FileID{
			Title: title,
			URL:   url,
		}
		// check if URL is in database before it even gets that far
		exists := a.checkFile(ctx, fileKey)
		// existing file
		if exists {
			results.ExistingFilesCount++
			// msg := fmt.Sprintf("File %s at %s already exists in database", fileKey.Title, fileKey.URL)
			// results.FailedFiles[fileKey] = msg
			continue
		} // if

		// if file is non duplicate, attempt to fetch
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			// log to FailedFiles
			msg := "non-200 status code"
			if err != nil {
				msg = err.Error()
			} // if
			results.FailedFiles[fileKey] = msg
			continue
		} // if
		// log to DownloadCount and SuccessFiles
		results.DownloadCount++
		results.SuccessFiles[fileKey] = resp.Body
		// resp.Body.Close()
	} // for
	if len(urlMap) > 0 && results.DownloadCount == 0 {
		err := errors.New("could not download any files")
		return results, err
	} // if
	return results, err
} // DownloadFiles
