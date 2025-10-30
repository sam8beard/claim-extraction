package shared

import (
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type File struct {
	FileName  string // name of file
	ObjectKey string // object key of file
	URL       string // url of file
	Status    string // status of file (downloaded, uploaded, failed)
}

type FailedFile struct {
	URL    string // url of file
	Report string // reason for why file failed in acquisition
}

type FileID struct {
	Title     string
	URL       string
	ObjectKey string
}

type MinioClient struct {
	Bucket string
	Client *minio.Client
}

type Workflow struct {
	MinioClient MinioClient   // bucket name with minio client
	PGClient    *pgxpool.Pool // pgx pool
}

type UploadResult struct {
	SuccessUploadCount int
	ExistingFilesCount int
	SuccessFiles       []File
	FailedFiles        []FailedFile
}

type DownloadResult struct {
	FailedFiles        map[FileID]string
	SuccessFiles       map[FileID]io.ReadCloser
	DownloadCount      int
	ExistingFilesCount int
}
