# Storage Package

The storage package provides an abstraction layer for file storage operations with support for multiple backends (Local filesystem and AWS S3).

## Features

- ✅ Unified interface for various storage backends
- ✅ Local filesystem storage
- ✅ AWS S3 compatible storage
- ✅ Stream support for large files
- ✅ Automatic directory creation (local)
- ✅ Public/private file access control
- ✅ File operations: Put, Get, Delete, Exists, GetURL
- ✅ Context support for timeout and cancellation

## Installation

```bash
go get github.com/budimanlai/go-pkg/storage
go get github.com/aws/aws-sdk-go-v2/service/s3
```

## Storage Interface

All storage backends implement the same interface:

```go
type Storage interface {
    Put(ctx context.Context, path string, file io.Reader, options *PutOptions) (string, error)
    Get(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
    Exists(ctx context.Context, path string) (bool, error)
    GetURL(ctx context.Context, path string) (string, error)
}
```

## Local Storage

### Basic Usage

```go
package main

import (
    "context"
    "os"
    
    "github.com/budimanlai/go-pkg/storage"
)

func main() {
    // Create local storage instance
    localStorage, err := storage.NewLocalStorage(storage.LocalStorageConfig{
        BasePath: "./uploads",
        BaseURL:  "http://localhost:3000/uploads",
    })
    if err != nil {
        panic(err)
    }
    
    // Upload file
    file, _ := os.Open("document.pdf")
    defer file.Close()
    
    url, err := localStorage.Put(context.Background(), "documents/doc.pdf", file, nil)
    if err != nil {
        panic(err)
    }
    
    println("File uploaded:", url)
    // Output: http://localhost:3000/uploads/documents/doc.pdf
}
```

### Configuration

```go
type LocalStorageConfig struct {
    // BasePath is the base directory for file storage
    BasePath string
    
    // BaseURL is the base URL for accessing files
    BaseURL string
}
```

### File Operations

**Upload File:**
```go
file, _ := os.Open("image.jpg")
defer file.Close()

url, err := localStorage.Put(context.Background(), "images/photo.jpg", file, nil)
```

**Download File:**
```go
reader, err := localStorage.Get(context.Background(), "images/photo.jpg")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Save to file
output, _ := os.Create("downloaded.jpg")
defer output.Close()
io.Copy(output, reader)
```

**Check if File Exists:**
```go
exists, err := localStorage.Exists(context.Background(), "images/photo.jpg")
if exists {
    println("File exists")
}
```

**Delete File:**
```go
err := localStorage.Delete(context.Background(), "images/photo.jpg")
```

**Get File URL:**
```go
url, err := localStorage.GetURL(context.Background(), "images/photo.jpg")
// Returns: http://localhost:3000/uploads/images/photo.jpg
```

## S3 Storage

### Basic Usage

```go
package main

import (
    "context"
    "os"
    
    "github.com/budimanlai/go-pkg/storage"
)

func main() {
    // Create S3 storage instance
    s3Storage, err := storage.NewS3Storage(storage.S3StorageConfig{
        Region:      "us-east-1",
        Bucket:      "my-bucket",
        AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
        SecretKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
    })
    if err != nil {
        panic(err)
    }
    
    // Upload file
    file, _ := os.Open("document.pdf")
    defer file.Close()
    
    url, err := s3Storage.Put(context.Background(), "documents/doc.pdf", file, &storage.PutOptions{
        ContentType: "application/pdf",
        Public:      true,
    })
    if err != nil {
        panic(err)
    }
    
    println("File uploaded:", url)
}
```

### Configuration

```go
type S3StorageConfig struct {
    // Region is the AWS region
    Region string
    
    // Bucket is the S3 bucket name
    Bucket string
    
    // Endpoint is optional for S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
    // Leave empty for AWS S3
    Endpoint string
    
    // AccessKeyID is the AWS access key
    AccessKeyID string
    
    // SecretKey is the AWS secret key
    SecretKey string
    
    // UsePathStyle forces path-style addressing (for MinIO and some S3-compatible services)
    UsePathStyle bool
}
```

### S3-Compatible Services

**MinIO:**
```go
s3Storage, err := storage.NewS3Storage(storage.S3StorageConfig{
    Region:       "us-east-1",
    Bucket:       "my-bucket",
    Endpoint:     "http://localhost:9000",
    AccessKeyID:  "minioadmin",
    SecretKey:    "minioadmin",
    UsePathStyle: true,
})
```

**DigitalOcean Spaces:**
```go
s3Storage, err := storage.NewS3Storage(storage.S3StorageConfig{
    Region:      "nyc3",
    Bucket:      "my-space",
    Endpoint:    "https://nyc3.digitaloceanspaces.com",
    AccessKeyID: "your-access-key",
    SecretKey:   "your-secret-key",
})
```

### Put Options

```go
type PutOptions struct {
    // ContentType specifies the MIME type of the file
    ContentType string
    
    // Public makes the file publicly accessible
    Public bool
}
```

**Upload with Options:**
```go
file, _ := os.Open("image.jpg")
defer file.Close()

url, err := s3Storage.Put(context.Background(), "images/photo.jpg", file, &storage.PutOptions{
    ContentType: "image/jpeg",
    Public:      true,
})
```

## Integration with Fiber

### File Upload Handler

```go
package main

import (
    "github.com/budimanlai/go-pkg/storage"
    "github.com/gofiber/fiber/v2"
)

var storageService storage.Storage

func main() {
    // Initialize storage
    storageService, _ = storage.NewLocalStorage(storage.LocalStorageConfig{
        BasePath: "./uploads",
        BaseURL:  "http://localhost:3000/uploads",
    })
    
    app := fiber.New()
    
    app.Post("/upload", uploadHandler)
    app.Get("/files/:path", downloadHandler)
    app.Delete("/files/:path", deleteHandler)
    
    app.Listen(":3000")
}

func uploadHandler(c *fiber.Ctx) error {
    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }
    
    // Open file
    src, err := file.Open()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to open file",
        })
    }
    defer src.Close()
    
    // Generate path
    path := fmt.Sprintf("uploads/%s", file.Filename)
    
    // Upload to storage
    url, err := storageService.Put(c.Context(), path, src, &storage.PutOptions{
        ContentType: file.Header.Get("Content-Type"),
        Public:      true,
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to upload file",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "url":     url,
    })
}

func downloadHandler(c *fiber.Ctx) error {
    path := c.Params("path")
    
    // Get file from storage
    reader, err := storageService.Get(c.Context(), path)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "File not found",
        })
    }
    defer reader.Close()
    
    // Stream file to response
    return c.SendStream(reader)
}

func deleteHandler(c *fiber.Ctx) error {
    path := c.Params("path")
    
    // Delete file
    err := storageService.Delete(c.Context(), path)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to delete file",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "File deleted",
    })
}
```

### Multiple File Upload

```go
func uploadMultipleHandler(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid form data",
        })
    }
    
    files := form.File["files"]
    urls := make([]string, 0, len(files))
    
    for _, file := range files {
        src, err := file.Open()
        if err != nil {
            continue
        }
        defer src.Close()
        
        path := fmt.Sprintf("uploads/%s", file.Filename)
        url, err := storageService.Put(c.Context(), path, src, nil)
        if err != nil {
            continue
        }
        
        urls = append(urls, url)
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "urls":    urls,
        "count":   len(urls),
    })
}
```

## Advanced Usage

### Factory Pattern for Multiple Backends

```go
package main

import (
    "context"
    "io"
    "os"
    
    "github.com/budimanlai/go-pkg/storage"
)

type StorageService struct {
    storage storage.Storage
}

func NewStorageService() (*StorageService, error) {
    var store storage.Storage
    var err error
    
    // Choose backend based on environment
    if os.Getenv("STORAGE_TYPE") == "s3" {
        store, err = storage.NewS3Storage(storage.S3StorageConfig{
            Region:      os.Getenv("AWS_REGION"),
            Bucket:      os.Getenv("S3_BUCKET"),
            AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
            SecretKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
        })
    } else {
        store, err = storage.NewLocalStorage(storage.LocalStorageConfig{
            BasePath: "./uploads",
            BaseURL:  "http://localhost:3000/uploads",
        })
    }
    
    if err != nil {
        return nil, err
    }
    
    return &StorageService{storage: store}, nil
}

func (s *StorageService) UploadFile(ctx context.Context, path string, file io.Reader) (string, error) {
    return s.storage.Put(ctx, path, file, nil)
}

func (s *StorageService) DownloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
    return s.storage.Get(ctx, path)
}

func (s *StorageService) DeleteFile(ctx context.Context, path string) error {
    return s.storage.Delete(ctx, path)
}
```

### Context Timeout

```go
import "time"

// Upload with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

url, err := storageService.Put(ctx, "large-file.zip", file, nil)
```

### Error Handling

```go
url, err := storageService.Put(ctx, "file.txt", file, nil)
if err != nil {
    switch {
    case errors.Is(err, storage.ErrNotFound):
        log.Println("File not found")
    case errors.Is(err, storage.ErrPermission):
        log.Println("Permission denied")
    case errors.Is(err, context.DeadlineExceeded):
        log.Println("Upload timeout")
    default:
        log.Println("Upload error:", err)
    }
}
```

## Best Practices

1. **Use Context for Timeout**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

2. **Close Readers**
   ```go
   reader, err := storageService.Get(ctx, path)
   if err != nil {
       return err
   }
   defer reader.Close()
   ```

3. **Validate File Types**
   ```go
   allowedTypes := map[string]bool{
       "image/jpeg": true,
       "image/png":  true,
   }
   
   if !allowedTypes[file.Header.Get("Content-Type")] {
       return errors.New("invalid file type")
   }
   ```

4. **Use Environment Variables for Credentials**
   ```go
   config := storage.S3StorageConfig{
       AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
       SecretKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
   }
   ```

5. **Generate Unique Filenames**
   ```go
   import "github.com/google/uuid"
   
   ext := filepath.Ext(file.Filename)
   filename := uuid.New().String() + ext
   path := fmt.Sprintf("uploads/%s", filename)
   ```

## Testing

The storage interface makes testing easy with mocks:

```go
type MockStorage struct {
    PutFunc    func(ctx context.Context, path string, file io.Reader, options *storage.PutOptions) (string, error)
    GetFunc    func(ctx context.Context, path string) (io.ReadCloser, error)
    DeleteFunc func(ctx context.Context, path string) error
}

func (m *MockStorage) Put(ctx context.Context, path string, file io.Reader, options *storage.PutOptions) (string, error) {
    return m.PutFunc(ctx, path, file, options)
}

// Implement other methods...

// Usage in tests
mockStorage := &MockStorage{
    PutFunc: func(ctx context.Context, path string, file io.Reader, options *storage.PutOptions) (string, error) {
        return "http://example.com/file.jpg", nil
    },
}
```

## Performance Tips

1. **Stream Large Files** - Jangan load seluruh file ke memory
2. **Use Context Timeout** - Prevent hanging requests
3. **Implement Retry Logic** - For transient network errors
4. **Cache File URLs** - If using S3 presigned URLs
5. **Use CDN** - For frequently accessed files

## Security Considerations

1. **Validate File Extensions**
2. **Check File Size Limits**
3. **Scan for Malware** (for user uploads)
4. **Use Presigned URLs** for temporary access
5. **Set Proper ACLs** on S3 buckets
6. **Never Commit Credentials** to repository

## See Also

- [Helpers Package](./helpers.md)
- [Response Package](./response/README.md)
