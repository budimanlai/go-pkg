package storage

type BaseStorage interface {
	// Save uploads a file from sourceFile path to the destination path in the storage system.
	Save(sourceFile string, destination string) error

	// Delete removes the file at the specified path from the storage system.
	Delete(path string) error

	// Exists checks if a file exists at the specified path in the storage system.
	Exists(path string) (bool, error)

	// GetURL generates a publicly accessible URL for the file at the specified path.
	GetURL(path string) (string, error)

	// GetSignedURL generates a signed URL for the file at the specified path with an expiry time in seconds.
	GetSignedURL(path string, expirySeconds int64) (string, error)
}
