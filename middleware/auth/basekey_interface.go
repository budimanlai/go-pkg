package auth

// BaseKey defines the interface for managing base keys in the authentication system.
type BaseKey interface {
	// Add adds a new base key to the system.
	Add(key string) error

	// AddKeyValue adds a new base key with an associated value to the system.
	// if the key already exists, it updates the value.
	AddKeyValue(key string, value string) error

	// Replace delete the existing base keys and adds new base keys with the associated values.
	Replace(keys map[string]string) error

	// Remove deletes an existing base key from the system.
	Remove(key string) error

	// RemoveAll deletes all base keys from the system.
	RemoveAll() error

	// GetValue retrieves the value associated with the given base key.
	GetValue(key string) (string, error)

	// IsExists checks if the base key exists and returns its identifier.
	IsExists(key string) bool
}
