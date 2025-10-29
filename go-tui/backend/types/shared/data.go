package shared

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type File struct {
	FileName string // name of file
	Key      string // object key of file
	URL      string // url of file
	Status   string // status of file (downloaded, uploaded, failed)
}

type FailedFile struct {
	URL    string // url of file
	Report string // reason for why file failed in acquisition
}

type MinioClient struct {
	Bucket string
	Client *minio.Client
}

type Workflow struct {
	MinioClient MinioClient   // bucket name with minio client
	PGClient    *pgxpool.Pool // pgx pool
}
