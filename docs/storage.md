# Storage Package

The storage package provides a unified abstraction for file storage operations, with built-in backends for the local filesystem and S3-compatible object storage (AWS S3, MinIO, SeaweedFS, DigitalOcean Spaces, etc.).

## Features

- ✅ Unified `BaseStorage` interface across backends
- ✅ Local filesystem storage
- ✅ S3-compatible storage (path-style addressing, custom endpoints)
- ✅ Upload from a file path or an `io.Reader`
- ✅ Automatic directory creation (local)
- ✅ Automatic MIME type detection on S3 uploads
- ✅ Public URL and time-limited signed URL generation
- ✅ File operations: Save, SaveFromReader, Delete, Exists, GetURL, GetSignedURL

## Installation

```bash
go get github.com/budimanlai/go-pkg/v3/storage
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/aws/aws-sdk-go-v2/feature/s3/manager
```

## Storage Interface

Both backends implement `BaseStorage`:

```go
type BaseStorage interface {
    // Save uploads a file from sourceFile path to the destination path in the storage system.
    Save(sourceFile string, destination string) error

    // SaveFromReader uploads a file from an io.Reader to the destination path in the storage system.
    SaveFromReader(reader io.Reader, destination string) error

    // Delete removes the file at the specified path from the storage system.
    Delete(path string) error

    // Exists checks if a file exists at the specified path in the storage system.
    Exists(path string) (bool, error)

    // GetURL generates a publicly accessible URL for the file at the specified path.
    GetURL(path string) (string, error)

    // GetSignedURL generates a signed URL for the file at the specified path with an expiry time in seconds.
    GetSignedURL(path string, expirySeconds int64) (string, error)
}
```

The `Storage` wrapper delegates every call to whichever backend it holds, so application code can depend on `*storage.Storage` without caring which backend is active:

```go
type Storage struct {
    Storage BaseStorage
}

func NewStorage(base BaseStorage) *Storage
```

## Local Storage

### Basic Usage

```go
package main

import (
    "fmt"

    "github.com/budimanlai/go-pkg/v3/storage"
)

func main() {
    localBackend := storage.NewLocalStorage("./uploads", "http://localhost:3000/uploads")
    fileStorage := storage.NewStorage(localBackend)

    err := fileStorage.Save("./data/document.pdf", "documents/doc.pdf")
    if err != nil {
        panic(err)
    }

    url, _ := fileStorage.GetURL("documents/doc.pdf")
    fmt.Println("File URL:", url)
    // Output: http://localhost:3000/uploads/documents/doc.pdf
}
```

### Constructor

```go
func NewLocalStorage(uploadDir, baseURL string) BaseStorage
```

- `uploadDir` — base directory on disk where files are written. Subdirectories are created automatically.
- `baseURL` — base URL prefixed to a path when building the value returned by `GetURL`.

`GetSignedURL` on local storage does not generate a real signature — it just returns the same value as `GetURL`, since local files have no expiring-access concept.

### File Operations

**Upload from a file path:**
```go
err := fileStorage.Save("./data/image.jpg", "images/photo.jpg")
```

**Upload from an `io.Reader`:**
```go
f, _ := os.Open("image.jpg")
defer f.Close()

err := fileStorage.SaveFromReader(f, "images/photo.jpg")
```

**Check if a file exists:**
```go
exists, err := fileStorage.Exists("images/photo.jpg")
```

**Delete a file:**
```go
err := fileStorage.Delete("images/photo.jpg")
```

**Get a public URL:**
```go
url, err := fileStorage.GetURL("images/photo.jpg")
// Returns: http://localhost:3000/uploads/images/photo.jpg
```

## S3 Storage

### Basic Usage

```go
package main

import (
    "fmt"

    "github.com/budimanlai/go-pkg/v3/storage"
)

func main() {
    config := storage.S3Config{
        Region:          "your_bucket_region",
        Bucket:          "your_bucket_name",
        AccessKeyID:     "your_access_key_id",
        SecretAccessKey: "your_secret_access_key",
        EndpointURL:     "your_endpoint_url",
    }

    s3Backend := storage.NewS3Storage(config)
    fileStorage := storage.NewStorage(s3Backend)

    err := fileStorage.Save("./data/image1.png", "public/avatar/image1.png")
    if err != nil {
        panic(err)
    }

    url, _ := fileStorage.GetURL("public/avatar/image1.png")
    fmt.Println("File URL:", url)

    signedURL, _ := fileStorage.GetSignedURL("public/avatar/image1.png", 60)
    fmt.Println("Signed URL (60s):", signedURL)
}
```

### Configuration

```go
type S3Config struct {
    // Region is the S3 region
    Region string

    // Bucket is the S3 bucket name
    Bucket string

    // AccessKeyID is the S3 access key
    AccessKeyID string

    // SecretAccessKey is the S3 secret key
    SecretAccessKey string

    // EndpointURL is the S3-compatible service endpoint (MinIO, SeaweedFS,
    // DigitalOcean Spaces, etc). Leave empty to use the default AWS S3 endpoint.
    // GetURL and GetSignedURL both build their URLs from this value, since the
    // client always uses path-style addressing (EndpointURL/Bucket/Key).
    EndpointURL string
}
```

`NewS3Storage` always enables path-style addressing (`UsePathStyle = true`) for compatibility with MinIO/SeaweedFS-style services. If `EndpointURL` is set, it becomes the S3 client's base endpoint; if left empty, the AWS SDK's default endpoint resolution for `Region` is used instead.

### S3-Compatible Services

**MinIO:**
```go
config := storage.S3Config{
    Region:          "us-east-1",
    Bucket:          "my-bucket",
    AccessKeyID:     "minioadmin",
    SecretAccessKey: "minioadmin",
    EndpointURL:     "http://localhost:9000",
}
```

**DigitalOcean Spaces:**
```go
config := storage.S3Config{
    Region:          "nyc3",
    Bucket:          "my-space",
    AccessKeyID:     "your-access-key",
    SecretAccessKey: "your-secret-key",
    EndpointURL:     "https://nyc3.digitaloceanspaces.com",
}
```

### Content Type Detection

`Save` and `SaveFromReader` automatically sniff the MIME type from the first 512 bytes of the uploaded content (via `http.DetectContentType`) and send it as the object's `ContentType` — there's no need to set it manually.

### File Operations

**Upload from a file path:**
```go
err := fileStorage.Save("./data/document.pdf", "documents/doc.pdf")
```

**Upload from an `io.Reader`:**
```go
f, _ := os.Open("image.jpg")
defer f.Close()

err := fileStorage.SaveFromReader(f, "images/photo.jpg")
```

**Check if a file exists:**
```go
exists, err := fileStorage.Exists("images/photo.jpg")
```

**Delete a file:**
```go
err := fileStorage.Delete("images/photo.jpg")
```

**Get a public URL:**
```go
url, err := fileStorage.GetURL("images/photo.jpg")
// Built as: EndpointURL/Bucket/images/photo.jpg
```

**Get a signed (time-limited) URL:**
```go
signedURL, err := fileStorage.GetSignedURL("images/photo.jpg", 300) // expires in 300 seconds
```

## Choosing a Backend at Runtime

```go
package main

import (
    "os"

    "github.com/budimanlai/go-pkg/v3/storage"
)

func newFileStorage() *storage.Storage {
    if os.Getenv("STORAGE_TYPE") == "s3" {
        backend := storage.NewS3Storage(storage.S3Config{
            Region:          os.Getenv("AWS_REGION"),
            Bucket:          os.Getenv("S3_BUCKET"),
            AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
            SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
            EndpointURL:     os.Getenv("S3_ENDPOINT_URL"),
        })
        return storage.NewStorage(backend)
    }

    backend := storage.NewLocalStorage("./uploads", "http://localhost:3000/uploads")
    return storage.NewStorage(backend)
}
```

## Best Practices

1. **Use environment variables for credentials**
   ```go
   config := storage.S3Config{
       AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
       SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
   }
   ```

2. **Generate unique filenames** to avoid collisions
   ```go
   import "github.com/google/uuid"

   ext := filepath.Ext(originalFilename)
   dest := fmt.Sprintf("uploads/%s%s", uuid.New().String(), ext)
   ```

3. **Prefer `SaveFromReader`** when the source is already in memory or streamed (e.g. an HTTP upload body) instead of writing it to a temp file first.

4. **Use signed URLs** (`GetSignedURL`) for private files that should only be accessible for a limited time, and `GetURL` only for objects that are actually public.

5. **Never commit credentials** to the repository — pass them via environment variables or a secrets manager.

## See Also

- [Helpers Package](./helpers.md)
- [Response Package](./response/README.md)
