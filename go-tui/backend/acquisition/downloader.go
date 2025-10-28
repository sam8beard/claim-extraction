package acquisition

import (
	"io"
	"net/http"
)

type FileKey struct {
	Title string
	URL   string
}

type DownloadResults struct {
	FailedFiles   map[FileKey]string
	SuccessFiles  map[FileKey]io.ReadCloser
	DownloadCount int
}

func DownloadFiles(urlMap map[string]string) DownloadResults {
	results := DownloadResults{
		FailedFiles:  make(map[FileKey]string),
		SuccessFiles: make(map[FileKey]io.ReadCloser, 0),
	}
	for title, url := range urlMap {
		fileKey := FileKey{
			Title: title,
			URL:   url,
		}
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			// log to FailedFiles
			results.FailedFiles[fileKey] = err.Error()
			continue
		} // if
		// log to DownloadCount and SuccessFiles
		results.DownloadCount++
		results.SuccessFiles[fileKey] = resp.Body
		resp.Body.Close()
	} // for
	return results
} // DownloadFiles
