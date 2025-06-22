package s3 

import ( 
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/aws/aws-sdk-go-v2/config"
	"context"
	"log"
)

func NewClient() (*s3.Client, error) { 
	err := godotenv.Load("../../.env")
	if err != nil { 
		log.Fatal("Failed to load .env file:", err)
		return nil, err
	} // if 

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil { 
		log.Fatal(err)
		return nil, err
	} // if 

	client := s3.NewFromConfig(cfg)
	
	return client, nil
} // client 