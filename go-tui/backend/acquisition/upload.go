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
	"tui/backend/db"
	"tui/backend/db/models"
	"tui/backend/types"

	"github.com/minio/minio-go/v7"
)

type UploadResults struct {
	SuccessUploadCount int
	ExistingFilesCount int
	SuccessFiles       []types.File
	FailedFiles        []types.FailedFile
}

func NewObjectKey(fileTitle string, fileURL string) (string, error) {
	hostName, err := NewHostName(fileURL)
	// construct object key
	fileTitle = fmt.Sprint(fileTitle, ".pdf")
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

func (a *Acquisition) Upload(files *map[FileKey]io.ReadCloser) (UploadResults, error) {
	var err error
	uploadResults := UploadResults{
		SuccessFiles: make([]types.File, 0),
		FailedFiles:  make([]types.FailedFile, 0),
	}

	for fileKey, fileReader := range *files {

		// unpack file data
		title, url := fileKey.Title, fileKey.URL

		// create copy of fileReader
		reader := fileReader

		// create temp file and new reader for seek and metadeta
		tempFile, err := os.CreateTemp("", "tempfile-*")
		if err != nil {
			err := errors.New("could not create temp file")
			return uploadResults, err
		} // if

		// copy file contents into tempFile and get size
		fileSize, err := io.Copy(tempFile, reader)
		if err != nil {
			err = errors.New("could not copy file contents")
			return uploadResults, err
		} // if

		// reset file offset
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			tempFile.Close()
			err = errors.New("could not reset file offset")
			return uploadResults, err
		} // if

		reader.Close()
		defer tempFile.Close()

		// get file hash
		data, err := io.ReadAll(tempFile)
		if err != nil {
			err = errors.New("could not read file body")
			return uploadResults, err
		}
		h := sha256.New()
		h.Write(data)
		hashBytes := h.Sum(nil)
		fileHash := hex.EncodeToString(hashBytes)

		// reset pointer
		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			tempFile.Close()
			err = errors.New("could not reset file offset")
			return uploadResults, err
		} // if

		// create object key for file upload
		hostName, err := NewHostName(url)
		if err != nil {
			return uploadResults, err
		} // if
		fileKey, err := NewObjectKey(title, url)
		if err != nil {
			return uploadResults, err
		} // if

		ctx := context.Background()
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
			FileName:      title,
			Source:        hostName,
			ContentHash:   fileHash,
			S3Key:         fileKey,
			FileSizeBytes: int(fileSize),
			TextExtracted: false,
		}

		// THIS IS WHERE WE DETERMINE WHAT FILES ARE ALREADY IN OUR DATABASE
		// WHAT SHOULD WE DO?
		// insert row
		err = db.InsertDocumentMetadata(ctx, a.PGClient, &doc)
		if err != nil {
			uploadResults.ExistingFilesCount++
			file := types.FailedFile{
				URL:    url,
				Report: "file already exists in database",
			}
			uploadResults.FailedFiles = append(uploadResults.FailedFiles, file)
		} else {
			uploadResults.SuccessUploadCount++
			file := types.File{
				FileName: title,
				Key:      fileKey,
				URL:      url,
				Status:   "downloaded",
			}
			uploadResults.SuccessFiles = append(uploadResults.SuccessFiles, file)
		} // if
	} // for
	return uploadResults, err
} // Upload
