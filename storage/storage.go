package storage

import "io"

type Storage struct {
	Storage BaseStorage
}

func NewStorage(base BaseStorage) *Storage {
	storage := &Storage{
		Storage: base,
	}

	return storage
}

// Save uploads a file from sourceFile path to the destination path in the storage system.
func (s *Storage) Save(sourceFile string, destination string) error {
	return s.Storage.Save(sourceFile, destination)
}

func (s *Storage) SaveFromReader(reader io.Reader, destination string) error {
	return s.Storage.SaveFromReader(reader, destination)
}

// Delete removes the file at the specified path from the storage system.
func (s *Storage) Delete(path string) error {
	return s.Storage.Delete(path)
}

// Exists checks if a file exists at the specified path in the storage system.
func (s *Storage) Exists(path string) (bool, error) {
	return s.Storage.Exists(path)
}

// GetURL generates a publicly accessible URL for the file at the specified path.
func (s *Storage) GetURL(path string) (string, error) {
	return s.Storage.GetURL(path)
}

func (s *Storage) GetSignedURL(path string, expirySeconds int64) (string, error) {
	return s.Storage.GetSignedURL(path, expirySeconds)
}
