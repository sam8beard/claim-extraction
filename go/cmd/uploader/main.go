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
	"log"
	"github.com/sam8beard/claim-extraction/go/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() { 
	var filePath string
	var source string 

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

	fmt.Println(filePath, source)

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

	// make file key for upload 
	fileName := filepath.Base(filePath)
	currTime := time.Now() 
	formattedTime := currTime.Format(time.RFC3339)
	fileKey := fmt.Sprint(source, "/", formattedTime, "_", fileName)
	log.Println(fileKey)
	
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

	fmt.Println("Successfully uploaded to: ", result)

	// establish connection to pg db
	pool, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil { 
		fmt.Println("Unable to establish database connection")
		return
	} // if 
	defer pool.Close()

} // main