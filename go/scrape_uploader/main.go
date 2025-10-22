package main
// package main


import ( 
	"fmt"
	"github.com/sam8beard/claim-extraction/go/s3client"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"context"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"io"
	"github.com/sam8beard/claim-extraction/go/utils"
	"github.com/sam8beard/claim-extraction/go/db"
	"github.com/sam8beard/claim-extraction/go/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"crypto/sha256"
	"encoding/hex"
	"github.com/sam8beard/claim-extraction/go/sqsclient"
	"net/url"
	// "math/big"
	// "bytes"
	// "github.com/gocolly/colly"
	// "net/http"
	// "slices"

	// "fmt"
	// "github.com/gocolly/colly"
	// "encoding/json"
	// "strings"
	// "slices"
	// "regexp"
	// "io"
	// "github.com/sam8beard/claim-extraction/go/scrape_upload"
)

func main() { 
	filesUploaded := 0
	// set env vars for db conn ection
	err := utils.LoadDotEnvUpwards()
	if err != nil { 
		fmt.Println("Could not load .env variables")
		fmt.Println(err)
		return
	} // if 

	// initialize all connections needed 
	// s3, sqs, sns, etc. 
	// connect to s3 client
	clientS3, err := s3client.NewClient() 
	if err != nil { 
		panic(err)
	} // if
	
	// connect to sqs client
	clientSQS, err := sqsclient.NewClient()
	if err != nil { 
		panic(err)
	} // if 
	
	queueUrl := "https://sqs.us-east-2.amazonaws.com/728951503252/claim-extraction-message-queue"

	// link s3 client to uploader
	uploader := manager.NewUploader(clientS3)
	
	// establish connection pool to pg db
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil { 
			fmt.Println("Unable to establish database connection")
			panic(err)
	} // if 
	defer pool.Close()

	// allFiles = map[io.ReadCloser]string[fileTitle, fileUrl]
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
			// panic(err)
			continue
		} // if 
		hostName := parsedUrl.Hostname()
		
		// construct object key
		fileTitle = fmt.Sprint(fileTitle, ".pdf")
		currTime := time.Now()
		formattedTime := currTime.Format(time.RFC3339)
		fileKey := fmt.Sprint("raw", "/", hostName, "/", formattedTime, "-", fileTitle)
		
		// insert into bucket 
		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("claim-pipeline-docstore"),
			Key: 	aws.String(fileKey),
			Body:	tempFile,
		})

		if err != nil { 
			fmt.Println("Could not upload to S3 bucket.")
			// panic(err)
			continue
		} // if 

		fmt.Println("Successfully uploaded to S3 bucket:  ", result, "\n")
		
		// wait for message to confirm file has been processed 
		output, err := clientSQS.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(queueUrl),
			WaitTimeSeconds: 20,
		})
		
		if err != nil { 
			fmt.Println("Could not receive message")
			panic(err)
		} // if 
		
		
		// create document entry
		doc := models.Document{ 
			FileName: fileTitle, 
			Source: hostName, 
			ContentHash: fileHash, 
			S3Key: fileKey, 
			FileSizeBytes: int(fileSize),
			TextExtracted: false, 
		}

		// If file was properly processed, then insert row with TextExtracted = true 
		// else, TextExtracted = false
		if len(output.Messages) > 0 { 
			doc.TextExtracted = true
		} // if

		fmt.Printf("Prepared document for insertion: %+v\n\n", doc)
	
		// insert row
		err = db.InsertDocumentMetadata(context.Background(), pool, &doc)
		if err != nil {
			fmt.Println("Error inserting row into database: ", err)
			continue
		} // if 
		
		fmt.Println("Successfully inserted row into Postgres database")
		filesUploaded += 1


		// construct file key
		//		get title from file map
		//		get base url 
		// 		key = raw, /, base url, time, -, title
	
		
		// upload to bucket with file reader and key

		// check for message 

		// establish connection pool to db 

		// create documents entry 

		// insert doc if processed 
		
	} // for 
	
	fmt.Printf("Files successfully uploaded: %d\n", filesUploaded)
	
	// link s3 uploader 


} // main 