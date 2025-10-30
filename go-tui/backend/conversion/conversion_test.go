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

type Lockbox struct {
	mu sync.Mutex
	r  Result
}
type Result struct {
	data map[shared.FileID][]byte
	err  map[shared.FileID]string
}

func (l *Lockbox) modifyLocker(f shared.FileID, u any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch val := u.(type) {
	case string:
		l.r.err[f] = val
	case []byte:
		l.r.data[f] = val
	}
}

func TestExtractTwo(t *testing.T) {
	msg := []byte("hello")
	fmsg := "hello"

	file := shared.FileID{
		Title: string("cool"),
	}
	fFile := shared.FileID{
		Title: string("not cool"),
	}

	l := Lockbox{
		r: Result{
			data: make(map[shared.FileID][]byte),
			err:  make(map[shared.FileID]string),
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	doSomethingGood := func() {
		defer wg.Done()
		l.modifyLocker(file, msg)
	}
	go doSomethingGood()
	doSomethingBad := func() {
		defer wg.Done()
		l.modifyLocker(fFile, fmsg)
	}
	go doSomethingBad()

	wg.Wait()

	t.Log("firing")
	t.Log(l.r.data)
	t.Log(l.r.err)
}
