package conversion

import (
	"context"
	"log"
	"testing"
	"tui/backend/acquisition"
	"tui/backend/types"
)

var ctx = context.Background()

func NewConversion() *Conversion {
	c := Conversion{}
	err := c.InitializeClients()
	if err != nil {
		log.Fatalf("failed to initialize clients")
	}
	return &c
} // NewAcquisition

func NewConversionInput() types.ConversionInput {
	a := acquisition.Acquisition{}
	a.InitializeClients()
	input := types.AcquisitionInput{
		Query:     "consulting white papers in 2025",
		FileCount: 50,
	}
	acResult, err := a.Run(ctx, input)
	if err != nil {
		log.Fatal(err)
	} // if

	conInput := types.ConversionInput{
		SuccessFiles:  acResult.SuccessFiles,
		FailedFiles:   acResult.FailedFiles,
		Log:           acResult.Log,
		PagesSearched: acResult.PagesSearched,
		URLsScraped:   acResult.URLsScraped,
	}
	return conInput
}

func TestDownload(t *testing.T) {
	input := NewConversionInput()
	con := NewConversion()

	_, err := con.Download(ctx, input)
	if err != nil {
		panic(err)
	}

}

func TestExtract(t *testing.T) {

	input := NewConversionInput()
	con := NewConversion()

	downloadResult, err := con.Download(ctx, input)
	if err != nil {
		panic(err)
	}
	_, err = con.Extract(ctx, downloadResult)
	if err != nil {
		t.Fatalf("Extraction failed")
	}
}

// func TestDownload(t *testing.T) {
// 	conv := NewAcquisition()
// 	input := types.ConversionInput{
// 		SuccessFiles: []shared.File{
// 			{
// 				Title:     "test.pdf",
// 				ObjectKey: "raw/test.pdf",
// 				URL:       "http://example.com/test.pdf",
// 			},
// 		},
// 	}
// 	// You may need to mock MinioClient here for isolated testing
// 	result, err := conv.Download(context.Background(), input)
// 	if err != nil {
// 		t.Fatalf("Download failed: %v", err)
// 	}
// 	if len(result.SuccessFiles) == 0 {
// 		t.Error("No files downloaded")
// 	}
// }

// func TestExtract(t *testing.T) {
// 	conv := &Conversion{}

// 	// Prepare a dummy DownloadResult with a sample PDF file
// 	file, err := os.Open("testdata/sample.pdf")
// 	if err != nil {
// 		t.Skip("sample.pdf not found, skipping Extract test")
// 	}
// 	defer file.Close()
// 	fileID := shared.FileID{
// 		Title:     "sample.pdf",
// 		ObjectKey: "raw/sample.pdf",
// 		URL:       "http://example.com/sample.pdf",
// 	}
// 	downloadResult := &shared.DownloadResult{
// 		SuccessFiles: map[shared.FileID]io.ReadCloser{
// 			fileID: file,
// 		},
// 		FailedFiles: make(map[shared.FileID]string),
// 	}
// 	result, err := conv.Extract(context.Background(), downloadResult)
// 	if err != nil {
// 		t.Fatalf("Extract failed: %v", err)
// 	}
// 	if len(result.SuccessFiles) == 0 {
// 		t.Error("No files extracted")
// 	}
// }

// func TestRun(t *testing.T) {
// 	conv := &conversion.Conversion{}

// 	input := types.ConversionInput{
// 		SuccessFiles: []shared.File{

// 			Title:     "test.pdf",
// 			ObjectKey: "raw/test.pdf",
// 			URL:       "http://example.com/test.pdf",
// 		},
// 	}
// 	// You may need to mock clients for isolated testing
// 	result, err := conv.Run(context.Background(), input)
// 	if err != nil {
// 		t.Fatalf("Run failed: %v", err)
// 	}
// 	if len(result.ConvertedFiles) == 0 {
// 		t.Error("No files converted")
// 	}
// }
