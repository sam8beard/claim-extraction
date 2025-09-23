/*
	Function to connect to SNS client
*/
package sqsclient

import ( 
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/config"
	"context"
	"log"
)

func NewClient() (*sqs.Client, error) { 

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil { 
		log.Fatal(err)
		return nil, err
	} // if 

	client := sqs.NewFromConfig(cfg)
	return client, nil

} // NewClient


