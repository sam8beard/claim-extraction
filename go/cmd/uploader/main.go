/*
	Entry point for file upload CLI tool

	Compilation: 
	cd go/cmd/uploader 
	./build.sh

	Usage: 
	./upload --file [path-to-file] --source [source-name]

	* NOTE * 
	Source names cannot be whitespace separated
	
*/
package main

import ( 
	"fmt"
	"flag"
	"github.com/sam8beard/claim-extraction/go/s3client"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
	"path/filepath"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"io"
	// "log"
	"github.com/sam8beard/claim-extraction/go/utils"
	"github.com/sam8beard/claim-extraction/go/db"
	"github.com/sam8beard/claim-extraction/go/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"crypto/sha256"
	"encoding/hex"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() { 
	var filePath string
	var source string 
	// var doc models.Document

	// set env vars 
	err := utils.LoadDotEnvUpwards()
	if err != nil { 
		fmt.Println("Could not load .env variables")
		return
	} // if 

	// identify which flags to look for - user must enter file path and source
	flag.StringVar(&filePath, "file", "", "Path to the document")
	flag.StringVar(&source, "source", "", "Source of the document")

	// parse flags provided 
	flag.Parse()

	if filePath == "" || source == "" { 
		fmt.Println("Must provide file path and source to use")
		return
	} // if 

	// connect to client
	client, err := s3client.NewClient() 
	if err != nil { 
		panic(err)
	} // if
	
	// get file reader for upload
	fileReader, err := os.Open(filePath)
	if err != nil { 
		panic(err)
	} // if 
	defer fileReader.Close() 
	
	// get file size for metadata 
	fileInfo, err := fileReader.Stat() 
	if err != nil { 
		fmt.Println(err)
		return
	} // if 
	fileSize := fileInfo.Size()

	// make new file reader for getting file hash for metadata
	fileReader2, err := os.Open(filePath)
	if err != nil { 
		panic(err)
	} // if 
	defer fileReader2.Close() 
	fileContents, err := io.ReadAll(fileReader2)
	if err != nil { 
		fmt.Println("Error reading file body")
		return
	} // if 
	h := sha256.New()
	h.Write(fileContents)
	hashBytes := h.Sum(nil)
	fileHash := hex.EncodeToString(hashBytes)

	// make file key for upload/metadata
	fileName := filepath.Base(filePath)
	currTime := time.Now() 
	formattedTime := currTime.Format(time.RFC3339)
	fileKey := fmt.Sprint("raw","/", source, "/", formattedTime, "-", fileName)
	
	// link client to uploader
	uploader := manager.NewUploader(client)

	// insert row in bucket
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("claim-pipeline-docstore"),
		Key: 	aws.String(fileKey),
		Body:	fileReader,
	})

	if err != nil { 
		panic(err)
	} // if 

	fmt.Println("Successfully uploaded to S3 bucket: ", result, "\n")
	
	
	
	// establish connection pool to pg db
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil { 
		fmt.Println("Unable to establish database connection")
		return
	} // if 
	defer pool.Close()
	
	// If file was properly processed, then insert row with TextExtracted = true 
	// else, TextExtracted = false
	doc := models.Document{ 
		FileName: fileName, 
		Source: source, 
		ContentHash: fileHash, 
		S3Key: fileKey, 
		FileSizeBytes: int(fileSize),
		TextExtracted: false, 
	} 
	
	fmt.Printf("Prepared document for insertion: %+v\n\n", doc)
	
	// insert row
	err = db.InsertDocumentMetadata(context.Background(), pool, &doc)
	if err != nil {
		fmt.Println("Error inserting row into database: ", err)
		return
	} // if 
	
	fmt.Println("Successfully inserted row into Postgres database")

	
	return 
} // main