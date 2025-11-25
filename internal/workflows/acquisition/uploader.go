package acquisition

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/db"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/db/models"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

type UploadResult struct {
	SuccessUploadCount int
	ExistingFilesCount int
	SuccessFiles       []shared.File
	FailedFiles        []shared.FailedFile
}

func NewObjectKey(fileTitle string, fileURL string) (string, error) {
	hostName, err := NewHostName(fileURL)
	// construct object key
	currTime := time.Now()
	formattedTime := currTime.Format(time.RFC3339)
	fileKey := fmt.Sprint("raw", "/", hostName, "/", formattedTime, "-", fileTitle)

	return fileKey, err
} // NewObjectKey

func NewHostName(fileURL string) (string, error) {
	// get host for object key
	parsedUrl, err := url.Parse(fileURL)
	if err != nil {
		err := errors.New("could not parse URL")
		return "", err
	} // if
	hostName := parsedUrl.Hostname()
	return hostName, err
}

func (a *Acquisition) Upload(ctx context.Context, files *map[shared.FileID]io.ReadCloser) (UploadResult, error) {
	var err error
	uploadResults := UploadResult{
		SuccessFiles: make([]shared.File, 0),
		FailedFiles:  make([]shared.FailedFile, 0),
	}

	for fileKey, fileReader := range *files {

		// unpack file data
		title, url := fileKey.Title, fileKey.URL

		// create copy of fileReader
		reader := fileReader

		// create temp file and new reader for seek and metadeta
		tempFile, err := os.CreateTemp("", "tempfile-*")
		if err != nil {
			file := shared.FailedFile{
				URL:    url,
				Report: "could not open temp file",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
			continue
		} // if

		// copy file contents into tempFile and get size
		fileSize, err := io.Copy(tempFile, reader)
		if err != nil {
			file := shared.FailedFile{
				URL:    url,
				Report: "could not copy file contents",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
			continue
		} // if

		// reset file offset
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			tempFile.Close()
			file := shared.FailedFile{
				URL:    url,
				Report: "could not reset offset",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
			continue

		} // if

		reader.Close()
		defer tempFile.Close()

		// get file hash
		data, err := io.ReadAll(tempFile)
		if err != nil {
			file := shared.FailedFile{
				URL:    url,
				Report: "could not read file body",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
			continue
		}
		h := sha256.New()
		h.Write(data)
		hashBytes := h.Sum(nil)
		fileHash := hex.EncodeToString(hashBytes)

		// reset pointer
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			tempFile.Close()
			file := shared.FailedFile{
				URL:    url,
				Report: "could not reset offset",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
			continue
		} // if

		// create object key for file upload
		hostName, err := NewHostName(url)
		if err != nil {
			return uploadResults, err
		} // if

		finalTitle := fmt.Sprint(title, ".pdf")
		fileKey, err := NewObjectKey(finalTitle, url)
		if err != nil {
			return uploadResults, err
		} // if

		opts := minio.PutObjectOptions{
			ContentType: "application/pdf",
		}
		bucketName := a.MinioClient.Bucket
		_, err = a.MinioClient.Client.PutObject(
			ctx,
			bucketName,
			fileKey,
			tempFile,
			int64(fileSize),
			opts,
		)
		if err != nil {
			err = errors.New("could not upload file to bucket")
			return uploadResults, err
		} // if

		// create document entry
		doc := models.Document{
			FileName:      finalTitle,
			Source:        hostName,
			ContentHash:   fileHash,
			S3Key:         fileKey,
			FileSizeBytes: int(fileSize),
			TextExtracted: false,
		}

		// insert row
		err = db.InsertDocumentMetadata(ctx, a.PGClient, &doc)
		if err != nil {
			uploadResults.ExistingFilesCount++
			file := shared.FailedFile{
				URL:    url,
				Report: "file already exists in database",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
		} else {
			uploadResults.SuccessUploadCount++
			file := shared.File{
				FileName:  title,
				ObjectKey: fileKey,
				URL:       url,
				Status:    "downloaded",
			}
			uploadResults.SuccessFiles = append(uploadResults.SuccessFiles, file)
		} // if
	} // for
	return uploadResults, err
} // Upload
