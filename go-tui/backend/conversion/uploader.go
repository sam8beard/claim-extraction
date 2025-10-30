package conversion

import (
	"context"
	"tui/backend/types/shared"
)

func (c *Conversion) Upload(ctx context.Context, e *ExtractionResult) (*shared.UploadResult, error) {
	var err error
	uploadResult := shared.UploadResult{
		SuccessFiles: make([]shared.File, 0),
		FailedFiles:  make([]shared.FailedFile, 0),
	}

	return &uploadResult, err
}
