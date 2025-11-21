package conversion

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

type UploadResult struct {
	SuccessFiles []shared.File
	FailedFiles  []shared.FailedFile
}

func (c *Conversion) Upload(ctx context.Context, files map[shared.FileID][]byte) (*UploadResult, error) {
	log.Println("FIRING IN conversion/upload.go")
	var err error
	uploadResult := UploadResult{
		SuccessFiles: make([]shared.File, 0),
		FailedFiles:  make([]shared.FailedFile, 0),
	}
	log.Printf("Starting upload...")
	log.Printf("Length of files map: %d", len(files))
	for fileID, body := range files {
		log.Println("FIRING INSIDE UPLOAD LOOP")
		log.Printf("file id original key: %s", fileID.OriginalKey)
		fileReader := bytes.NewReader(body)

		fileSize := len(body)
		opts := minio.PutObjectOptions{
			ContentType: "text",
		}
		log.Printf("Putting object %s...", fileID.ObjectKey)
		_, err := c.MinioClient.Client.PutObject(
			ctx,
			c.MinioClient.Bucket,
			fileID.ObjectKey,
			fileReader,
			int64(fileSize),
			opts,
		)

		// could not upload file
		if err != nil {
			msg := "failed to upload text file to MinIO"
			fFile := shared.FailedFile{

				URL:    fileID.URL,
				Report: msg,
			}
			uploadResult.FailedFiles = append(uploadResult.FailedFiles, fFile)
			log.Printf("Upload failed")
			continue
		} // if
		log.Printf("Upload successful")

		log.Printf("Updating row...")
		// could not update row of file
		log.Printf("Original key, should be raw/.../.pdf: %s", fileID.OriginalKey)
		if err := c.Update(ctx, fileID); err != nil {
			log.Println("conversion/uploader.go FIRING ON Update ERROR")
			msg := fmt.Sprintf("unable to update row in documents: %s", fileID.ObjectKey)
			fFile := shared.FailedFile{

				URL:    fileID.URL,
				Report: msg,
			}
			uploadResult.FailedFiles = append(uploadResult.FailedFiles, fFile)
			log.Printf("Update failed")
			continue
		} // if

		// add successfully uploaded and updated file to our upload result
		sFile := shared.File{
			ObjectKey: fileID.ObjectKey,
			URL:       fileID.URL,
			FileName:  fileID.Title,
			Status:    "processed",
		}
		uploadResult.SuccessFiles = append(uploadResult.SuccessFiles, sFile)
	} // for

	if len(files) > 0 && len(uploadResult.SuccessFiles) == 0 {
		log.Print("No files to upload")
		return nil, errors.New("failed to upload any extracted files")
	}
	return &uploadResult, err
}
