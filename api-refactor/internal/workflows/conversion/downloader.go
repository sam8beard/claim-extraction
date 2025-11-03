package conversion

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

func (c *Conversion) Download(ctx context.Context, input types.ConversionInput) (*shared.DownloadResult, error) {
	var err error
	result := shared.DownloadResult{
		SuccessFiles: make(map[shared.FileID]io.ReadCloser),
		FailedFiles:  make(map[shared.FileID]string),
	}

	bucket := c.MinioClient.Bucket
	filesToUpload := input.SuccessFiles

	for _, acquiredFile := range filesToUpload {

		title := acquiredFile.FileName
		objectKey := acquiredFile.ObjectKey
		url := acquiredFile.URL
		// dont think we need status here?
		fileID := shared.FileID{
			Title:     title,
			URL:       url,
			ObjectKey: objectKey,
		}
		// attempt to download file
		ctx := context.Background()
		opts := minio.GetObjectOptions{}
		object, err := c.MinioClient.Client.GetObject(
			ctx,
			bucket,
			objectKey,
			opts,
		)
		if err != nil {
			msg := fmt.Sprintf("failed to download pdf file \n%s\n%s:%s", title, url, err.Error())
			result.FailedFiles[fileID] = msg
			continue
		} // if

		result.DownloadCount++
		result.SuccessFiles[fileID] = object
	} // for

	// if no files were downloaded successfully
	if len(result.SuccessFiles) == 0 {
		err = errors.New("failed to download any files")
	} // if
	return &result, err
} // GetPDFs
