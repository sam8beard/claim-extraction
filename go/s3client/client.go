package s3client

import ( 
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go-v2/config"
	"context"
	"os"
	// "path/filepath"
	"log"
)

func NewClient() (*s3.Client, error) { 

	err := godotenv.Load("../../../.env")
	if err != nil { 
		log.Fatal("Failed to load .env file:", err)
		return nil, err
	} // if 

	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	os.Setenv("AWS_REGION", os.Getenv("AWS_REGION"))

	log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Println(os.Getenv("AWS_REGION"))
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil { 
		log.Fatal(err)
		return nil, err
	} // if 
	log.Println(cfg.Region)
	creds, err := cfg.Credentials.Retrieve(ctx)
	log.Println("AccessKeyID:", creds.AccessKeyID)
	log.Println("SecretAccessKey:", creds.SecretAccessKey)
	client := s3.NewFromConfig(cfg)
	// log.Println(client.Credentials)
	return client, nil
} // client 