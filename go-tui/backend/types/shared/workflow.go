package shared

import (
	"context"
	"errors"
	"os"
	"tui/backend/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IsWorkflow interface {
	InitializeClients() error
	Run() error
}

func (w *Workflow) InitializeClients() error {
	// set env vars for db conn ection
	err := utils.LoadDotEnvUpwards()
	if err != nil {
		err := errors.New("could not load .env variables")
		return err
	} // if

	// create MinIO client
	endpoint := "localhost:9000"
	accessKeyID := "muel"
	secretAccessKey := "password"
	useSSL := false
	bucketName := "claim-pipeline-docstore"

	w.MinioClient.Bucket = bucketName
	w.MinioClient.Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		err := errors.New("unable to establish connection to MinIO")
		return err
	} // if

	// establish connection pool to pg db
	w.PGClient, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		err := errors.New("unable to establish connection to Postgres")
		return err
	} // if

	return err
}
