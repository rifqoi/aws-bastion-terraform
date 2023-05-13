package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDB() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.Region = "us-east-1"
		return nil
	})

	if err != nil {
		log.Panic("cannot load config:", err)
	}

	db := dynamodb.NewFromConfig(cfg)

	return db
}
