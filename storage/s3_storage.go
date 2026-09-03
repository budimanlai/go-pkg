package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	EndpointURL     string
}

type S3Storage struct {
	Config        S3Config
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
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

		// Set custom endpoint if EndpointURL is provided
		if s3Config.EndpointURL != "" {
			o.BaseEndpoint = aws.String(s3Config.EndpointURL)
		}
	})

	// Create a presigner
	presigner := s3.NewPresignClient(client)
	uploader := manager.NewUploader(client)

	return &S3Storage{
		Config:        s3Config,
		client:        client,
		presignClient: presigner,
		uploader:      uploader,
	}
}

func (s3s *S3Storage) Save(sourceFile string, destination string) error {
	// Open the source file
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	return s3s.SaveFromReader(file, destination)
}

func (s3s *S3Storage) SaveFromReader(reader io.Reader, destination string) error {
	// Clean the destination path
	key := filepath.ToSlash(filepath.Clean(destination))
	key = strings.TrimPrefix(key, "/")

	// Detect the content type from the file content
	contentType, body, err := detectContentType(reader)
	if err != nil {
		return fmt.Errorf("failed to detect content type: %w", err)
	}

	// Upload the file to S3
	_, err = s3s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3s.Config.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}

// detectContentType sniffs the MIME type from the first bytes of reader,
// returning a new reader that still yields the full original content.
func detectContentType(reader io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)
	n, err := io.ReadFull(reader, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", nil, err
	}

	contentType := http.DetectContentType(buf[:n])
	return contentType, io.MultiReader(bytes.NewReader(buf[:n]), reader), nil
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

	// Combine endpoint URL with bucket and path (path-style addressing)
	urlStr := strings.TrimSuffix(s3s.Config.EndpointURL, "/") + "/" + s3s.Config.Bucket
	if cleanPath != "" {
		urlStr += "/" + cleanPath
	}

	return urlStr, nil
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
