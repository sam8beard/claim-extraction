package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/gocolly/colly"
)

type ResponseBody struct {
	Results []Result `json:"results"`
	// Query string								`json:"query"`
	// NumOfResults float64						`json:"number_of_results"`
	// ResultsList []map[string][]Result				`json:"results"`

	// Answers map[string]string				`json:"query"`
	// Corrections map[string]string			`json:"query"`
	// Suggestions map[string]string			`json:"query"`
	// UnresponsiveEngines map[string]string	`json:"query"`

} // ResponseBody

type Result struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	// Content string			`json:"content"`
	// Thumbnail string		`json:"thumbnail"`
	// Engine string			`json:"engine"`
	// Template string			`json:"template"`
	// ParsedUrls []string		`json:"parsed_urls"`
	// ImgSrc string			`json:"img_src"`
	// Priority string			`json:"priority"`
	// Engines []string		`json:"engines"`
	// Positions []float64		`json:"positions"`
	// Score float64			`json:"score"`
	// Category string			`json:"category"`
	PublishedDate string `json:"publishedDate"`
} // Result

func GetFiles() map[io.ReadCloser][]string {
	// all unique urls retrieved
	var fileUrls []string
	// file title mapped to url
	files := make(map[string]string)
	// files that were correctly downloaded
	// the file reader mapped to the title, url
	finalFiles := make(map[io.ReadCloser][]string)

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {

		// get response body
		if r.StatusCode == 200 {
			body := r.Body

			// use this if we want to also scrape for files that are not located directly at url
			// headers := r.Headers

			// unmarshal response body
			var response ResponseBody
			json.Unmarshal(body, &response)
			// add urls
			for _, result := range response.Results {
				// check if url already exists in our slice
				// fmt.Println(result.Title)
				// fmt.Println(result.URL)

				if !slices.Contains(fileUrls, result.URL) {
					fileUrls = append(fileUrls, result.URL)
					lower := strings.ToLower(result.Title)
					noElp := strings.ReplaceAll(lower, "...", "")
					fileName := strings.ReplaceAll(noElp, " ", "-")
					reNonAlnum := regexp.MustCompile(`[^a-z0-9\-]`)
					fileName = reNonAlnum.ReplaceAllString(fileName, "-")
					reMultDash := regexp.MustCompile(`-+`)
					fileName = reMultDash.ReplaceAllString(fileName, "-")
					fileName = strings.Trim(fileName, "-")
					files[fileName] = result.URL
				} // if
			} // for

		} // if

	}) // OnResponse

	// base url
	url := "http://localhost:8888/search?q=AI+safety+ethics+filetype:pdf&format=json&pageno="

	// adjust for number of desired pages worth of requests
	pageLimit := 15
	for pageNum := 1; pageNum < pageLimit; pageNum++ {

		newUrl := fmt.Sprintf("%s%d", url, pageNum)
		c.Visit(newUrl)
	} // for

	// download files
	for title, url := range files {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("\nCould not retrieve file at %s\n", url)
			// resp.Close = true
			// panic(err)
		} else {
			body := resp.Body
			fmt.Printf("\nFile successfully retrieved: %s\n", url)
			finalFiles[body] = []string{title, url}
			// fmt.Println(body)
		} // if
		// resp.Close = true
	} // for

	// // Testing
	// for _, url := range fileUrls {
	// 	fmt.Println(url)
	// } // for

	// Testing
	// for reader, details := range finalFiles {
	// 	fmt.Println(reader, " : ", details[0], " : ", details[1])
	// } // for

	return finalFiles
} // GetFiles
