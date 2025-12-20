package auth

import (
	"errors"
	"sync"
)

// BaseKeyProvider is a concrete implementation of the BaseKey interface.
type BaseKeyProvider struct {
	keys map[string]string
	mu   sync.RWMutex
}

// NewBaseKeyProvider creates a new instance of BaseKeyProvider.
func NewBaseKeyProvider() BaseKey {
	return &BaseKeyProvider{
		keys: make(map[string]string),
	}
}

// Add adds a new base key to the system.
func (b *BaseKeyProvider) Add(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.keys[key] = key
	return nil
}

// AddKeyValue adds a new base key with an associated value to the system.
// if the key already exists, it updates the value.
func (b *BaseKeyProvider) AddKeyValue(key string, value string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.keys[key] = value

	return nil
}

// Replace delete the existing base keys and adds new base keys with the associated values.
func (b *BaseKeyProvider) Replace(keys map[string]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Clear existing keys
	b.keys = make(map[string]string)

	// Add new keys
	for k, v := range keys {
		b.keys[k] = v
	}
	return nil
}

// Remove deletes an existing base key from the system.
func (b *BaseKeyProvider) Remove(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.keys, key)
	return nil
}

// RemoveAll deletes all base keys from the system.
func (b *BaseKeyProvider) RemoveAll() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.keys = make(map[string]string)
	return nil
}

// GetValue retrieves the value associated with the given base key.
func (b *BaseKeyProvider) GetValue(key string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	value, exists := b.keys[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

// IsExists checks if the base key exists and returns its identifier.
func (b *BaseKeyProvider) IsExists(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if _, exists := b.keys[key]; exists {
		return true
	}
	return false
}
