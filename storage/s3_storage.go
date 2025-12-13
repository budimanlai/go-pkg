package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	ServerURL       string
	BaseURL         string
}

type S3Storage struct {
	Config        S3Config
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewS3Storage(s3Config S3Config) BaseStorage {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Config.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s3Config.AccessKeyID,
			s3Config.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config: %v", err))
	}

	// Create S3 client with custom options
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Set path-style addressing for compatibility with MinIO, SeaweedFS, etc.
		o.UsePathStyle = true

		// Set custom endpoint if BaseURL is provided
		if s3Config.ServerURL != "" {
			o.BaseEndpoint = aws.String(s3Config.ServerURL)
		}
	})

	// Create a presigner
	presigner := s3.NewPresignClient(client)

	return &S3Storage{
		Config:        s3Config,
		client:        client,
		presignClient: presigner,
	}
}

func (s3s *S3Storage) Save(sourceFile string, destination string) error {
	// Open the source file
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Clean the destination path
	key := filepath.ToSlash(filepath.Clean(destination))
	key = strings.TrimPrefix(key, "/")

	// Upload the file to S3
	_, err = s3s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3s.Config.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

func (s3s *S3Storage) Delete(path string) error {
	// Clean the path
	key := filepath.ToSlash(filepath.Clean(path))
	key = strings.TrimPrefix(key, "/")

	// Delete the file from S3
	_, err := s3s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3s.Config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (s3s *S3Storage) Exists(path string) (bool, error) {
	// Clean the path
	key := filepath.ToSlash(filepath.Clean(path))
	key = strings.TrimPrefix(key, "/")

	// Check if the file exists in S3
	_, err := s3s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s3s.Config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check for NoSuchKey error
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, nil
		}
		// Check for NotFound error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to check file existence in S3: %w", err)
	}

	return true, nil
}

func (s3s *S3Storage) GetURL(path string) (string, error) {
	// Clean the path and replace backslashes with forward slashes for URLs
	cleanPath := filepath.ToSlash(filepath.Clean(path))

	// Remove leading slash if exists to avoid double slashes in URL
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	// Combine base URL with path
	url := s3s.Config.BaseURL
	if !strings.HasSuffix(url, "/") && cleanPath != "" {
		url += "/"
	}
	url += cleanPath

	return url, nil
}

func (s3s *S3Storage) GetSignedURL(path string, expirySeconds int64) (string, error) {
	// Clean the path
	key := filepath.ToSlash(filepath.Clean(path))
	key = strings.TrimPrefix(key, "/")

	// Generate the presigned URL
	presignedURL, err := s3s.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3s.Config.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires((time.Duration(expirySeconds) * time.Second)))
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return presignedURL.URL, nil
}
