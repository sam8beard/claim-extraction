package acquisition

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/sam8beard/claim-extraction/internal/types"
)

var ctx = context.Background()

func NewAcqInput() types.AcquisitionInput {
	input := types.AcquisitionInput{
		Query:     "consulting white papers in 2025",
		FileCount: 20,
	}
	return input
} // NewAcqInput

func NewAcquisition(t *testing.T) *Acquisition {
	a := Acquisition{}
	err := a.InitializeClients()
	if err != nil {
		t.Fatalf("failed to initialize clients")
	}
	return &a
} // NewAcquisition

// PASS
func TestAcquisition(t *testing.T) {
	file, err := os.Create("testing.log")
	if err != nil {
		t.Fatal(err)
	} // if
	log.SetOutput(file)
	a := NewAcquisition(t)
	input := NewAcqInput()
	ctx := context.Background()
	output, err := a.Run(ctx, input)
	if err != nil {
		t.Fatal(err)
	} // err
	resultLog := output.Log
	for _, result := range resultLog {
		log.Printf("\n%s\n", result)
	} // for

} // TestAcquisition

// PASS
func TestInitializeClients(t *testing.T) {
	a := Acquisition{}
	err := a.InitializeClients()
	if err != nil {
		t.Fatalf("failed to initialize clients")
	} // if
	t.Logf("%v", a.MinioClient)
	t.Logf("%v", a.PGClient)
} // TestInitializeClients

// PASS
func TestScrape(t *testing.T) {
	input := NewAcqInput()
	scrapeResult, err := Scrape(ctx, input.Query, input.FileCount)
	if err != nil {
		t.Fatal(err)
	} // if

	for fileName, fileURL := range scrapeResult.URLMap {
		t.Log(fileName, fileURL)
	} // for
	t.Log(scrapeResult.PageCount)
	t.Log(scrapeResult.URLCount)

} // TestScrape

// PASS
func TestDownload(t *testing.T) {
	a := NewAcquisition(t)
	input := NewAcqInput()
	scrapeResult, _ := Scrape(ctx, input.Query, input.FileCount)
	downloadResult, err := a.Download(ctx, scrapeResult.URLMap)
	if err != nil {
		t.Fatal(err)
	} // if

	for fileKey, report := range downloadResult.FailedFiles {
		// title := fileKey.Title
		url := fileKey.URL
		t.Logf("%s failed to download: %s", url, report)
	} // if

	for fileKey, _ := range downloadResult.SuccessFiles {
		url := fileKey.URL
		t.Logf("SUCCESS: %s", url)
	} // for

	t.Logf("Files successfully downloaded: %d", downloadResult.DownloadCount)
} // TestDownload

// PASS
func TestUpload(t *testing.T) {
	a := NewAcquisition(t)
	input := NewAcqInput()
	scrapeResult, _ := Scrape(ctx, input.Query, input.FileCount)
	downloadResult, err := a.Download(ctx, scrapeResult.URLMap)
	if err != nil {
		t.Fatal(err)
	} // if

	uploadResult, err := a.Upload(ctx, &downloadResult.SuccessFiles)
	if err != nil {
		t.Fatal(err)
	}

	for _, succFile := range uploadResult.SuccessFiles {
		t.Logf("%s\t%s\t%s\t%s", succFile.FileName, succFile.ObjectKey, succFile.URL, succFile.Status)
	} // for
	t.Log("\n\n\n")
	for _, failedFile := range uploadResult.FailedFiles {
		t.Logf("%s\t%s", failedFile.URL, failedFile.Report)
	} // for

	t.Logf("\nSuccesses: %d", uploadResult.SuccessUploadCount)
	t.Logf("\nExisting: %d", uploadResult.ExistingFilesCount)

} // TestUpload
