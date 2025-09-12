package providers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func CreateAWSConfig(ctx context.Context, endpoint, region string) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if endpoint != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region),
			config.WithBaseEndpoint(endpoint), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				"test",
				"test",
				"",
			)),
		)

	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}

	return cfg, err
}
