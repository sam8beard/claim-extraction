package main

import ( 
	"fmt"
	"github.com/gocolly/colly"
	"encoding/json"
	// "strings"
	// "io"
	// "slices"
)

// endpoints?:  
// http://localhost:8888/search?q=AI+safety+ethics&categories=files&format=json&pageno=1,2,3,4,...
// http://localhost:8888/search?q=AI+safety+ethics+filetype:pdf&format=json&pageno=1,2,3,4,...

// in order to get multiple pages, we are going to have to manually and iteratively change the 
// url extension pageno= value. the api does not support returning the number of pages of a query

type ResponseBody struct { 
	Results []Result	`json:"results"`
	// Query string								`json:"query"`
	// NumOfResults float64						`json:"number_of_results"`
	// ResultsList []map[string][]Result				`json:"results"`
	
	// Answers map[string]string				`json:"query"`
	// Corrections map[string]string			`json:"query"`
	// Suggestions map[string]string			`json:"query"`
	// UnresponsiveEngines map[string]string	`json:"query"`

} // ResponseBody 

type Result struct { 
	URL string				`json:"url"`
	// Title string			`json:"title"`
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
	// PublishedDate string	`json:"publishedDate"`
}

func main() { 
	var file_urls []string
	c := colly.NewCollector()

	// iteratively increase 1 up to maybeeee 15?
	// url := "http://localhost:8888/search?q=AI+safety+ethics+filetype:pdf&format=json&pageno=1"
	

	// for each page no
	//		send request to inital page load url
	// 		get the json response 
	//		add urls to file_url slice 
	// for each url in file_url slice
	//		if file url hasnt been seen
	//			download file from from unique url 
	// 			get file meta deta 
	// 			construct s3 entry and postgres row entry
	// 			upload to s3 bucket and insert row
	

	c.OnResponse(func(r *colly.Response) { 
		
		fmt.Println(r.StatusCode)
		// get response body 
		body := r.Body
		
		// unmarshal response body
		var response ResponseBody
		json.Unmarshal(body, &response)

		// add urls 
		for _, result := range response.Results {
			fmt.Println(result)
			// fmt.Println(url_string)
			// for i, url := range result.URL {
			// 	fmt.Println(url)
			// 	// check if url already exist in our list 
			// 	if !slices.Contains(file_urls, url) { 
			// 		// add url
			// 		file_urls = append(file_urls, url)
			// 	} // if 
			// }
		} // for 
	}) // OnResponse

	
	
	url := "http://localhost:8888/search?q=AI+safety+ethics+filetype:pdf&format=json&pageno="

	for page_num := 1; page_num < 16; page_num++ { 
		
		new_url := fmt.Sprintf("%s%d", url, page_num)
		c.Visit(new_url)
	} // for 

	for _, url := range file_urls { 
		fmt.Println(url)
	} // for 
	
} // main 