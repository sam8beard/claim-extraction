/*
Handles the orchestration of the acquisition flow
*/
package acquisition

import (
	"context"
	"errors"
	"fmt"
	"os"
	"tui/backend/types"
	"tui/backend/utils"

	// "github.com/minio/minio-go"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Bucket string
	Client *minio.Client
}

// represents an acquistion instance
type Acquisition struct {
	MinioClient MinioClient   // bucket name with minio client
	PGClient    *pgxpool.Pool // pgx pool
}

// // creates an acquisition instance
// func NewAcquisition(minioClient *minio.Client, pgClient *pgxpool.Pool) *Acquisition {
// 	// return &Acquisition{
// 	// 	MinioClient: minioClient,
// 	// 	PGClient:    pgClient,
// 	// 	// Scraper:     scraper,
// 	// }
// } // NewAcquisition

func (a *Acquisition) InitializeClients() error {
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

	a.MinioClient.Bucket = bucketName
	a.MinioClient.Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		err := errors.New("unable to establish connection to MinIO")
		return err
	} // if

	// establish connection pool to pg db
	a.PGClient, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		err := errors.New("unable to establish connection to Postgres")
		return err
	} // if
	defer a.PGClient.Close()

	return err
} // NewClients

// Executes the acquisition flow
func (a *Acquisition) Run(input types.AcquisitionInput) (types.AcquisitionResult, error) {
	var err error
	result := types.AcquisitionResult{
		SuccessFiles: make([]types.File, 0),
		FailedFiles:  make([]types.FailedFile, 0),
		Log:          make([]string, 0),
	}

	// 1) scrape urls
	scrapeResult := Scrape(input.Query, input.FileCount)

	// log count of pages searched to result and result log
	result.PagesSearched = scrapeResult.PageCount
	pagesSearchedMsg := fmt.Sprintf("%d out of a maximum 30 pages worth of results scraped", result.PagesSearched)
	result.Log = append(result.Log, pagesSearchedMsg)

	// log count of file urls scraped to result log
	result.URLsScraped = scrapeResult.URLCount
	urlCountMsg := fmt.Sprintf("%d out of a requested %d URLs scraped", scrapeResult.URLCount, input.FileCount)
	result.Log = append(result.Log, urlCountMsg)

	// 2) download files
	urlsToDownload := scrapeResult.URLMap
	downloadResults := DownloadFiles(urlsToDownload)

	// log results of downloaded files
	for fileInfo, _ := range downloadResults.SuccessFiles {
		title, url := fileInfo.Title, fileInfo.URL
		successMsg := fmt.Sprintf("File [%s] downloaded successfully from [%s]", title, url)
		result.Log = append(result.Log, successMsg)
	} // for

	// log total amount of downloaded files
	downloadedFileCount := len(downloadResults.SuccessFiles)
	successMsg := fmt.Sprintf("%d files downloaded successfully out of %d URLs", downloadedFileCount, result.URLsScraped)
	result.Log = append(result.Log, successMsg)

	// log results of failed files
	for fileInfo, report := range downloadResults.FailedFiles {
		title, url := fileInfo.Title, fileInfo.URL
		failedMsg := fmt.Sprintf("%s: Could not download %s from %s", report, title, url)
		result.Log = append(result.Log, failedMsg)
	} // for

	// log total amount of failed file
	failedFileCount := len(downloadResults.FailedFiles)
	failedMsg := fmt.Sprintf("Failed to download %d files", failedFileCount)
	result.Log = append(result.Log, failedMsg)

	// initialize clients
	if err := a.InitializeClients(); err != nil {
		return result, err
	} // if

	// 3) upload to minio
	if err := a.UploadToMinio(&downloadResults.SuccessFiles); err != nil {
		return result, err
	} // if
	// 4) save metadata to postgres
	// if err := a.UploadDataToPostgres()

	// 5) populate AcquisitionResult.Log

	// 6) return AcquisitionResult
	return result, err
} // Run
