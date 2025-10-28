package acquisition

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

func NewObjectKey(fileTitle string, fileURL string) (string, error) {
	// get host for object key
	parsedUrl, err := url.Parse(fileURL)
	if err != nil {
		err := errors.New("could not parse URL")
		return "", err
	} // if
	hostName := parsedUrl.Hostname()
	// construct object key
	fileTitle = fmt.Sprint(fileTitle, ".pdf")
	currTime := time.Now()
	formattedTime := currTime.Format(time.RFC3339)
	fileKey := fmt.Sprint("raw", "/", hostName, "/", formattedTime, "-", fileTitle)

	return fileKey, err
} // NewObjectKey

func (a *Acquisition) UploadToMinio(files *map[FileKey]io.ReadCloser) error {
	var err error
	for fileKey, fileReader := range *files {

		// unpack file data
		title, url := fileKey.Title, fileKey.URL

		// create copy of fileReader
		reader := fileReader

		// create temp file and new reader for seek and metadeta
		tempFile, err := os.CreateTemp("", "tempfile-*")
		if err != nil {
			err := errors.New("could not create temp file")
			return err
		} // if

		fileSize, err := io.Copy(tempFile, reader)
		if err != nil {
			err = errors.New("could not copy file contents")
			return err
		} // if
		fileKey, err := NewObjectKey(title, url)
		if err != nil {
			return err
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
		} // if
		return err
	} // for
	return err
} // UploadToMinio
