package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
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
	PublicURL       string
	PrivateURL      string
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

	// Upload the file to S3
	_, err := s3s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3s.Config.Bucket),
		Key:    aws.String(key),
		Body:   reader,
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
	urlStr := s3s.Config.PublicURL
	if !strings.HasSuffix(urlStr, "/") && cleanPath != "" {
		urlStr += "/"
	}
	urlStr += cleanPath

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

	// If PublicURL is provided, replace with custom domain and path
	if s3s.Config.PrivateURL != "" {
		parsedPresigned, err := url.Parse(presignedURL.URL)
		if err != nil {
			return "", fmt.Errorf("failed to parse presigned URL: %w", err)
		}

		// Parse the public URL
		parsedPublic, err := url.Parse(s3s.Config.PublicURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse public URL: %w", err)
		}

		// Build new URL with public domain and path
		// PublicURL might have path prefix like /storage/secure/
		publicPath := strings.TrimSuffix(parsedPublic.Path, "/")

		// Combine public path with file key
		finalPath := publicPath
		if key != "" {
			if publicPath != "" {
				finalPath = publicPath + "/" + key
			} else {
				finalPath = "/" + key
			}
		}

		// Build the final URL
		finalURL := url.URL{
			Scheme:   parsedPublic.Scheme,
			Host:     parsedPublic.Host,
			Path:     finalPath,
			RawQuery: parsedPresigned.RawQuery, // Keep signature and query params
		}

		return finalURL.String(), nil
	}

	return presignedURL.URL, nil
}
