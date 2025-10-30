package conversion

import (
	"context"
	"log"
	"sync"
	"testing"

	"tui/backend/acquisition"
	"tui/backend/types"
	"tui/backend/types/shared"
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

type testData struct {
	data     map[shared.FileID]string
	byteData map[shared.FileID][]byte
}
type testLocker struct {
	mu sync.Mutex
	td testData
}

func (l *testLocker) logEntry(id shared.FileID, u any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch val := u.(type) {
	case string:
		l.td.data[id] = val
	case []byte:
		l.td.byteData[id] = val
	} // switch
} // logEntry

func TestMutex(t *testing.T) {
	doSomething := func() {
		testString := "testing"
		fileToAdd := shared.FileID{
			ObjectKey: testString,
		}
		logEntry(fileToAdd)
	}
}
