package images

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/phzfi/RIC/server/config"
)

type S3Client struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3Client(cfg cfg.ImageSourceConfig) (*S3Client, error) {
	ctx := context.Background()

	var awsCfg aws.Config
	var err error

	if cfg.S3Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               cfg.S3Endpoint,
				HostnameImmutable: true,
			}, nil
		})

		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.S3Region),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.S3Region),
		)
	}

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		client: client,
		bucket: cfg.S3Bucket,
		prefix: cfg.S3Prefix,
	}, nil
}

func (s *S3Client) GetObject(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectKey := s.prefix + key

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("S3 GetObject error: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("S3 ReadAll error: %w", err)
	}

	return data, nil
}

func (s *S3Client) HeadObject(key string) (size int64, exists bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectKey := s.prefix + key

	_, err = s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return 0, false, nil
	}

	return 0, true, nil
}