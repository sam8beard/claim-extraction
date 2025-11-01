package conversion

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"tui/backend/types/shared"
)

func (c *Conversion) Upload(ctx context.Context, e *ExtractionResult) (*shared.UploadResult, error) {
	var err error
	uploadResult := shared.UploadResult{
		SuccessFiles: make([]shared.File, 0),
		FailedFiles:  make([]shared.FailedFile, 0),
	}

	extractedFiles := e.SuccessFiles

	for fileID, body := range extractedFiles {

		fileReader := bytes.NewReader(body)

		fileSize := len(body)
		opts := minio.PutObjectOptions{
			ContentType: "text",
		}
		_, err := c.MinioClient.Client.PutObject(
			ctx,
			c.MinioClient.Bucket,
			fileID.ObjectKey,
			fileReader,
			int64(fileSize),
			opts,
		)
		if err != nil {
			msg := "failed to text file to minio"
			fFile := shared.FailedFile{

				URL:    fileID.URL,
				Report: msg,
			}
			uploadResult.FailedFiles = append(uploadResult.FailedFiles, fFile)

		} // if

		// call Update() here
		sFile := shared.File{
			ObjectKey: fileID.ObjectKey,
			URL:       fileID.URL,
			FileName:  fileID.Title,
			Status:    "processed",
		}
		uploadResult.SuccessFiles = append(uploadResult.SuccessFiles, sFile)
	} // for

	return &uploadResult, err
}
