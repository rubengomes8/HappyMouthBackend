package kvstore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(ctx context.Context, url, region string) (*dynamodb.Client, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-east-1" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
