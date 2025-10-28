package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sam8beard/claim-extraction/go/db"
	"github.com/sam8beard/claim-extraction/go/models"
	"github.com/sam8beard/claim-extraction/go/utils"
)

func main() {
	// set env vars for db conn ection
	err := utils.LoadDotEnvUpwards()
	if err != nil {
		fmt.Println("Could not load .env variables")
		fmt.Println(err)
		return
	} // if

	// create MinIO client
	endpoint := "localhost:9000"
	accessKeyID := "muel"
	secretAccessKey := "password"
	useSSL := false
	bucketName := "claim-pipeline-docstore"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		fmt.Println("Unable to establish connection to MinIO server")
		panic(err)
	} // if

	// establish connection pool to pg db
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to establish database connection")
		panic(err)
	} // if
	defer pool.Close()

	// download files
	allFiles := GetFiles()

	// process all files
	for reader, details := range allFiles {
		// get file reader from reader slice
		fileReader := reader
		defer fileReader.Close()
		fileTitle := details[0]
		fileUrl := details[1]

		// create temp file and new reader for seek and metadeta
		tempFile, err := os.CreateTemp("", "tempfile-*")
		if err != nil {
			fmt.Println("Could not create temp file")
			// panic(err)
			continue
		} // if

		_, err = io.Copy(tempFile, fileReader)
		if err != nil {
			fmt.Printf("Unable to copy file for metadata: %s\n", fileTitle)
			continue
		} // if

		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			tempFile.Close()
			fmt.Println("Unable to reset offset")
			continue
		} // if

		fileReader.Close()
		defer tempFile.Close()

		fileInfo, err := tempFile.Stat()
		if err != nil {
			fmt.Println("Couldn't get file stats")
			// panic(err)
			continue
		} // if

		fileSize := fileInfo.Size()
		// fileSizeInt := big.NewInt(fileSize)
		if fileSize == 0 {
			fmt.Println("File is either empty or was not read properly")
			continue
		}

		// get file hash then reset pointer
		data, err := io.ReadAll(tempFile)
		if err != nil {
			fmt.Println("Error reading file body")
			// panic(err)
			continue
		} // if
		h := sha256.New()
		h.Write(data)
		hashBytes := h.Sum(nil)
		fileHash := hex.EncodeToString(hashBytes)
		tempFile.Seek(0, io.SeekStart)

		// get host for object key
		parsedUrl, err := url.Parse(fileUrl)
		if err != nil {
			fmt.Println("Could not parse url")
			continue
		} // if
		hostName := parsedUrl.Hostname()

		// construct object key
		fileTitle = fmt.Sprint(fileTitle, ".pdf")
		currTime := time.Now()
		formattedTime := currTime.Format(time.RFC3339)
		fileKey := fmt.Sprint("raw", "/", hostName, "/", formattedTime, "-", fileTitle)

		ctx := context.Background()
		opts := minio.PutObjectOptions{
			ContentType: "application/pdf",
		}
		_, err = minioClient.PutObject(ctx, bucketName, fileKey, tempFile, int64(fileSize), opts)
		if err != nil {
			fmt.Println("Could not upload file to bucket")
			continue
		} // if

		// fmt.Println(info)
		// create document entry
		doc := models.Document{
			FileName:      fileTitle,
			Source:        hostName,
			ContentHash:   fileHash,
			S3Key:         fileKey,
			FileSizeBytes: int(fileSize),
			TextExtracted: false,
		}

		// insert row
		err = db.InsertDocumentMetadata(context.Background(), pool, &doc)
		if err != nil {
			fmt.Println("File already exists in table: ", err)
			continue
		} // if

		fmt.Println("Successfully inserted row into Postgres database")
		// filesUploaded += 1
	} // for

} // main
