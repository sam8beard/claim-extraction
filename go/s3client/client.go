package s3client

import ( 
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"
	"context"
	"github.com/sam8beard/claim-extraction/go/utils"
	"log"
)

func NewClient() (*s3.Client, error) { 

	err := utils.LoadDotEnvUpwards()
	if err != nil { 
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