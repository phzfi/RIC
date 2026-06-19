package cache

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/phzfi/RIC/server/config"
)

type S3Store struct {
	sync.RWMutex
	client  *s3.Client
	bucket  string
	prefix  string
	maxSize uint64
}

func NewS3Store(cfg cfg.CacheConfig) (*S3Store, error) {
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

	maxSize := cfg.S3MaxMB * 1024 * 1024
	if maxSize == 0 {
		maxSize = 200 * 1024 * 1024
	}

	return &S3Store{
		client:  client,
		bucket:  cfg.S3Bucket,
		prefix:  cfg.S3Prefix,
		maxSize: maxSize,
	}, nil
}

func (s *S3Store) Load(key string) ([]byte, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectKey := s.prefix + stringToBase64(key)

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("S3 Load error for key %s: %v", objectKey, err)
		return nil, false
	}
	defer result.Body.Close()

	data := make([]byte, s.maxSize)
	n, err := result.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Printf("S3 Read error for key %s: %v", objectKey, err)
		return nil, false
	}

	return data[:n], true
}

func (s *S3Store) Store(key string, blob []byte) {
	if uint64(len(blob)) > s.maxSize {
		log.Printf("S3 Store: blob size %d exceeds max %d, skipping", len(blob), s.maxSize)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		objectKey := s.prefix + stringToBase64(key)

		_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(objectKey),
			Body:   bytes.NewReader(blob),
		})
		if err != nil {
			log.Printf("S3 Store error for key %s: %v", objectKey, err)
		}
	}()
}

func (s *S3Store) Delete(key string) uint64 {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		objectKey := s.prefix + stringToBase64(key)

		_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			log.Printf("S3 Delete error for key %s: %v", objectKey, err)
		}
	}()

	return 0
}