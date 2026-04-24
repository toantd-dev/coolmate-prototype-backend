package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coolmate/ecommerce-backend/internal/config"
)

type S3Manager struct {
	client *s3.Client
	bucket string
}

func NewS3Manager(cfg *config.S3Config) (*S3Manager, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint")
	})

	sdkConfig, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Manager{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (sm *S3Manager) Upload(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, src); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	key := fmt.Sprintf("%s/%d-%s", folder, time.Now().Unix(), file.Filename)

	_, err = sm.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(sm.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return key, nil
}

func (sm *S3Manager) Delete(ctx context.Context, key string) error {
	_, err := sm.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sm.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (sm *S3Manager) GetURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", "http://localhost:9000", sm.bucket, key)
}
