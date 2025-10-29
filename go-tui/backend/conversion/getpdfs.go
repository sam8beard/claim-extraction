package conversion

import (
	"io"

	"tui/backend/types/shared"
)

type PDFResult struct {
	PDFSuccess   map[shared.File]io.Reader // files and readers
	PDFFailed    []shared.FailedFile       // failed files
	SuccessCount int                       // count of successes
	FailedCount  int                       // count of failures
}

func (c *Conversion) GetPDFs() (*PDFResult, error) {
	var err error
	result := PDFResult{
		PDFSuccess: make(map[shared.File]io.Reader, 0),
		PDFFailed:  make([]shared.FailedFile, 0),
	}

	return &result, err
} // GetPDFs
