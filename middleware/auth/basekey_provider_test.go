package auth

import (
	"sync"
	"testing"
)

func TestNewBaseKeyImpl(t *testing.T) {
	impl := NewBaseKeyProvider()
	if impl == nil {
		t.Fatal("Expected NewBaseKeyImpl to return a non-nil instance")
	}

	// Check that it implements the interface
	var _ BaseKey = impl
}

func TestBaseKeyImpl_Add(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test adding a key
	err := impl.Add("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the key exists
	if !impl.IsExists("test-key") {
		t.Error("Expected key to exist after adding")
	}

	// Verify the value equals the key
	value, err := impl.GetValue("test-key")
	if err != nil {
		t.Errorf("Expected no error getting value, got %v", err)
	}
	if value != "test-key" {
		t.Errorf("Expected value to be 'test-key', got '%s'", value)
	}
}

func TestBaseKeyImpl_AddKeyValue(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test adding a key-value pair
	err := impl.AddKeyValue("key1", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the key exists
	if !impl.IsExists("key1") {
		t.Error("Expected key to exist after adding")
	}

	// Verify the value
	value, err := impl.GetValue("key1")
	if err != nil {
		t.Errorf("Expected no error getting value, got %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value to be 'value1', got '%s'", value)
	}

	// Test updating an existing key
	err = impl.AddKeyValue("key1", "value2")
	if err != nil {
		t.Errorf("Expected no error updating key, got %v", err)
	}

	value, err = impl.GetValue("key1")
	if err != nil {
		t.Errorf("Expected no error getting updated value, got %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected value to be 'value2', got '%s'", value)
	}
}

func TestBaseKeyImpl_Replace(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Add some initial keys
	impl.AddKeyValue("key1", "value1")
	impl.AddKeyValue("key2", "value2")

	// Replace with new keys
	newKeys := map[string]string{
		"key3": "value3",
		"key4": "value4",
	}

	err := impl.Replace(newKeys)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify old keys are removed
	if impl.IsExists("key1") {
		t.Error("Expected key1 to be removed")
	}
	if impl.IsExists("key2") {
		t.Error("Expected key2 to be removed")
	}

	// Verify new keys exist
	if !impl.IsExists("key3") {
		t.Error("Expected key3 to exist")
	}
	if !impl.IsExists("key4") {
		t.Error("Expected key4 to exist")
	}

	// Verify values
	value3, _ := impl.GetValue("key3")
	if value3 != "value3" {
		t.Errorf("Expected value3, got '%s'", value3)
	}

	value4, _ := impl.GetValue("key4")
	if value4 != "value4" {
		t.Errorf("Expected value4, got '%s'", value4)
	}
}

func TestBaseKeyImpl_Replace_EmptyMap(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Add some initial keys
	impl.AddKeyValue("key1", "value1")

	// Replace with empty map
	err := impl.Replace(map[string]string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all keys are removed
	if impl.IsExists("key1") {
		t.Error("Expected key1 to be removed")
	}
}

func TestBaseKeyImpl_Remove(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Add a key
	impl.AddKeyValue("key1", "value1")

	// Remove the key
	err := impl.Remove("key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the key is removed
	if impl.IsExists("key1") {
		t.Error("Expected key to be removed")
	}

	// Removing a non-existent key should not error
	err = impl.Remove("non-existent")
	if err != nil {
		t.Errorf("Expected no error removing non-existent key, got %v", err)
	}
}

func TestBaseKeyImpl_RemoveAll(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Add multiple keys
	impl.AddKeyValue("key1", "value1")
	impl.AddKeyValue("key2", "value2")
	impl.AddKeyValue("key3", "value3")

	// Remove all keys
	err := impl.RemoveAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all keys are removed
	if impl.IsExists("key1") {
		t.Error("Expected key1 to be removed")
	}
	if impl.IsExists("key2") {
		t.Error("Expected key2 to be removed")
	}
	if impl.IsExists("key3") {
		t.Error("Expected key3 to be removed")
	}
}

func TestBaseKeyImpl_GetValue(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test getting a non-existent key
	_, err := impl.GetValue("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
	if err.Error() != "key not found" {
		t.Errorf("Expected 'key not found' error, got '%s'", err.Error())
	}

	// Add a key and get its value
	impl.AddKeyValue("key1", "value1")
	value, err := impl.GetValue("key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%s'", value)
	}
}

func TestBaseKeyImpl_IsExists(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test non-existent key
	if impl.IsExists("non-existent") {
		t.Error("Expected key to not exist")
	}

	// Add a key and test
	impl.AddKeyValue("key1", "value1")
	if !impl.IsExists("key1") {
		t.Error("Expected key to exist")
	}

	// Remove and test again
	impl.Remove("key1")
	if impl.IsExists("key1") {
		t.Error("Expected key to not exist after removal")
	}
}

func TestBaseKeyImpl_ConcurrentAccess(t *testing.T) {
	impl := NewBaseKeyProvider()
	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			key := "key"
			value := "value"
			impl.AddKeyValue(key, value)
		}(i)
	}
	wg.Wait()

	// Concurrent reads and writes
	wg.Add(numGoroutines * 2)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			impl.AddKeyValue("concurrent-key", "concurrent-value")
		}(i)

		go func(idx int) {
			defer wg.Done()
			impl.GetValue("concurrent-key")
			impl.IsExists("concurrent-key")
		}(i)
	}
	wg.Wait()

	// Verify no race conditions occurred
	if !impl.IsExists("concurrent-key") {
		t.Error("Expected concurrent-key to exist")
	}
}

func TestBaseKeyImpl_MultipleOperations(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test a sequence of operations
	impl.Add("key1")
	impl.AddKeyValue("key2", "value2")
	impl.AddKeyValue("key3", "value3")

	if !impl.IsExists("key1") || !impl.IsExists("key2") || !impl.IsExists("key3") {
		t.Error("Expected all keys to exist")
	}

	impl.Remove("key1")
	if impl.IsExists("key1") {
		t.Error("Expected key1 to be removed")
	}

	impl.Replace(map[string]string{
		"key4": "value4",
		"key5": "value5",
	})

	if impl.IsExists("key2") || impl.IsExists("key3") {
		t.Error("Expected old keys to be removed after replace")
	}

	if !impl.IsExists("key4") || !impl.IsExists("key5") {
		t.Error("Expected new keys to exist after replace")
	}

	impl.RemoveAll()
	if impl.IsExists("key4") || impl.IsExists("key5") {
		t.Error("Expected all keys to be removed")
	}
}

func TestBaseKeyImpl_EmptyKeyAndValue(t *testing.T) {
	impl := NewBaseKeyProvider()

	// Test with empty key
	err := impl.Add("")
	if err != nil {
		t.Errorf("Expected no error with empty key, got %v", err)
	}

	if !impl.IsExists("") {
		t.Error("Expected empty key to exist")
	}

	// Test with empty key and value
	err = impl.AddKeyValue("", "")
	if err != nil {
		t.Errorf("Expected no error with empty key and value, got %v", err)
	}

	value, err := impl.GetValue("")
	if err != nil {
		t.Errorf("Expected no error getting empty key, got %v", err)
	}
	if value != "" {
		t.Errorf("Expected empty value, got '%s'", value)
	}
}

func TestBaseKeyImpl_SpecialCharacters(t *testing.T) {
	impl := NewBaseKeyProvider()

	specialKeys := []string{
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
		"key with spaces",
		"key@with#special$chars",
		"key/with/slashes",
		"UPPERCASE-KEY",
		"MixedCase-Key",
	}

	// Add all special keys
	for _, key := range specialKeys {
		err := impl.AddKeyValue(key, "value-"+key)
		if err != nil {
			t.Errorf("Failed to add key '%s': %v", key, err)
		}
	}

	// Verify all keys exist and have correct values
	for _, key := range specialKeys {
		if !impl.IsExists(key) {
			t.Errorf("Expected key '%s' to exist", key)
		}

		value, err := impl.GetValue(key)
		if err != nil {
			t.Errorf("Failed to get value for key '%s': %v", key, err)
		}
		expectedValue := "value-" + key
		if value != expectedValue {
			t.Errorf("Expected value '%s', got '%s'", expectedValue, value)
		}
	}
}
