package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	UploadDir string
	BaseURL   string
}

func NewLocalStorage(uploadDir, baseURL string) BaseStorage {
	return &LocalStorage{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}

func (ls *LocalStorage) Save(sourceFile string, destination string) error {
	// Construct the full destination path
	destPath := filepath.Join(ls.UploadDir, destination)

	// Create the directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open the source file
	srcFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy the file content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (ls *LocalStorage) Delete(path string) error {
	// Construct the full file path
	filePath := filepath.Join(ls.UploadDir, path)

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", err)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (ls *LocalStorage) Exists(path string) (bool, error) {
	// Construct the full file path
	filePath := filepath.Join(ls.UploadDir, path)

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

func (ls *LocalStorage) GetURL(path string) (string, error) {
	// Clean the path and replace backslashes with forward slashes for URLs
	cleanPath := filepath.ToSlash(filepath.Clean(path))

	// Remove leading slash if exists to avoid double slashes in URL
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	// Combine base URL with path using path.Join for URLs
	url := ls.BaseURL
	if !strings.HasSuffix(url, "/") && cleanPath != "" {
		url += "/"
	}
	url += cleanPath

	return url, nil
}

func (ls *LocalStorage) GetSignedURL(path string, expirySeconds int64) (string, error) {
	// For local storage, signed URLs are not typically implemented.
	// We will return the regular URL.
	return ls.GetURL(path)
}
