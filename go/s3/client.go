package s3 

import ( 
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"fmt"
	"context"
	"log"
)

func NewClient() (*s3.Client) { 
	err := godotenv.Load("../../.env")
	if err != nil { 
		log.Fatal("Failed to load .env file:", err)
	} // if 

	ctx := context.TODO()
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil { 
		log.Fatal(err)
	} // if 

	client := sts.NewFromConfig(cfg)

	return client
} // client 