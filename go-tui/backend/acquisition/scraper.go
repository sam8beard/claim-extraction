package acquisition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/gocolly/colly"
)

/*
Handles scraping for files

DOES NOT DOWNLOAD (see download.go)
*/

// the result of our scrape
type ScrapeResult struct {
	URLMap    map[string]string // map of file names to urls
	URLCount  int
	PageCount int
}

// json response body
type ResponseBody struct {
	Results []Result `json:"results"`
}

// content of json body
type Result struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// maximum pages allowed to visit
const MaxPages = 15 // 50

// builds the base url
func BuildBaseUrl(query string) string {
	baseUrl := "http://localhost:8888/search?q=%s+filetype:pdf&format=json&pageno="
	encodedQuery := strings.ReplaceAll(query, " ", "+")
	finalUrl := fmt.Sprintf(baseUrl, encodedQuery)
	return finalUrl
}

// builds and returns a cleaned file name
func BuildFileName(fileTitle string) string {
	lower := strings.ToLower(fileTitle)
	noElp := strings.ReplaceAll(lower, "...", "")
	fileName := strings.ReplaceAll(noElp, " ", "-")
	reNonAlnum := regexp.MustCompile(`[^a-z0-9\-]`)
	fileName = reNonAlnum.ReplaceAllString(fileName, "-")
	reMultDash := regexp.MustCompile(`-+`)
	fileName = reMultDash.ReplaceAllString(fileName, "-")
	fileName = strings.Trim(fileName, "-")
	return fileName
}

// makes a new query for each page iteration
func MakeQueryURL(query string, pageNo int) (string, error) {
	if pageNo >= MaxPages {
		err := errors.New("max page limit reached")
		return "", err
	} // if

	return fmt.Sprintf("%s%d", query, pageNo), nil
} // MakeQueryUrl

func ResponseHandler(ctx context.Context, r *colly.Response, finalResults *ScrapeResult, seenURLs *[]string, maxFiles int) {
	// invalid response
	if r.StatusCode != 200 {
		fmt.Println(r.Headers)
		return
	} // if
	body := r.Body
	var response ResponseBody
	err := json.Unmarshal(body, &response)
	// could not unmarshal
	if err != nil {
		return
	} // if

	for _, result := range response.Results {
		if !slices.Contains(*seenURLs, result.URL) && finalResults.URLCount < maxFiles {
			*seenURLs = append(*seenURLs, result.URL)
			fileName := BuildFileName(result.Title)
			finalResults.URLMap[fileName] = result.URL
			finalResults.URLCount++
			if finalResults.URLCount >= maxFiles {
				break
			} // if
		} // if
	} // for
} // ResponseHandler

/*
Scrapes for files using query and number of files provided by user

Returns a ScrapeResult
*/
func Scrape(ctx context.Context, query string, maxFiles int) (ScrapeResult, error) {
	seenURLs := make([]string, 0)
	var err error
	finalResults := ScrapeResult{
		URLMap: make(map[string]string),
	}
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		ResponseHandler(ctx, r, &finalResults, &seenURLs, maxFiles)
	})

	baseURL := BuildBaseUrl(query)
	pageNum := 1
	for {
		currURL, err := MakeQueryURL(baseURL, pageNum)
		// we have reached the maximum page limit
		if err != nil || finalResults.URLCount >= maxFiles {
			break
		} // if
		err = c.Visit(currURL)
		if err != nil {
			return finalResults, err
		} // if
		pageNum++
	} // for

	finalResults.PageCount = pageNum
	return finalResults, err
} // Scrape
