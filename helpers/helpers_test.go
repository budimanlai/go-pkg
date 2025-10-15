package helpers

import (
	"strings"
	"testing"
)

func TestPointer(t *testing.T) {
	// Test with int
	val := 42
	ptr := Pointer(val)
	if *ptr != val {
		t.Errorf("Expected %d, got %d", val, *ptr)
	}

	// Test with string
	str := "hello"
	strPtr := Pointer(str)
	if *strPtr != str {
		t.Errorf("Expected %s, got %s", str, *strPtr)
	}

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	person := Person{Name: "John", Age: 30}
	personPtr := Pointer(person)
	if *personPtr != person {
		t.Errorf("Expected %+v, got %+v", person, *personPtr)
	}
}

func TestDerefPointer(t *testing.T) {
	// Test with non-nil pointer
	val := 42
	ptr := &val
	result := DerefPointer(ptr, 0)
	if result != val {
		t.Errorf("Expected %d, got %d", val, result)
	}

	// Test with nil pointer
	var nilPtr *int
	result = DerefPointer(nilPtr, 100)
	if result != 100 {
		t.Errorf("Expected default value 100, got %d", result)
	}

	// Test with string
	str := "world"
	strPtr := &str
	resultStr := DerefPointer(strPtr, "default")
	if resultStr != str {
		t.Errorf("Expected %s, got %s", str, resultStr)
	}

	// Test with nil string pointer
	var nilStrPtr *string
	resultStr = DerefPointer(nilStrPtr, "fallback")
	if resultStr != "fallback" {
		t.Errorf("Expected fallback, got %s", resultStr)
	}
}

func TestUnmarshalTo(t *testing.T) {
	// Test successful unmarshal
	jsonStr := `{"name":"Alice","age":25}`
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person, err := UnmarshalTo[Person](jsonStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if person.Name != "Alice" {
		t.Errorf("Expected name Alice, got %s", person.Name)
	}
	if person.Age != 25 {
		t.Errorf("Expected age 25, got %d", person.Age)
	}

	// Test unmarshal to map
	jsonMapStr := `{"key":"value","number":123}`
	result, err := UnmarshalTo[map[string]interface{}](jsonMapStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal map: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("Expected key value, got %v", result["key"])
	}
	if result["number"] != float64(123) {
		t.Errorf("Expected number 123, got %v", result["number"])
	}

	// Test invalid JSON
	_, err = UnmarshalTo[Person]("invalid json")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestGenerateTrxID(t *testing.T) {
	id1 := GenerateTrxID()
	id2 := GenerateTrxID()

	// Check length (YYMMDDHHMMSS = 12 + 4 random = 16)
	if len(id1) != 16 {
		t.Errorf("Expected length 16, got %d", len(id1))
	}
	if len(id2) != 16 {
		t.Errorf("Expected length 16, got %d", len(id2))
	}

	// Check they are different (highly unlikely to be same due to timestamp + random)
	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}

	// Check format (should be numeric)
	for _, char := range id1 {
		if char < '0' || char > '9' {
			t.Errorf("ID should contain only digits, got %c", char)
		}
	}
}

func TestGenerateTrxIDWithPrefix(t *testing.T) {
	prefix := "TXN"
	id := GenerateTrxIDWithPrefix(prefix)

	if !strings.HasPrefix(id, prefix) {
		t.Errorf("Expected ID to start with %s, got %s", prefix, id)
	}

	// Total length should be prefix + 16
	expectedLen := len(prefix) + 16
	if len(id) != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, len(id))
	}
}

func TestGenerateTrxIDWithSuffix(t *testing.T) {
	suffix := "END"
	id := GenerateTrxIDWithSuffix(suffix)

	if !strings.HasSuffix(id, suffix) {
		t.Errorf("Expected ID to end with %s, got %s", suffix, id)
	}

	// Total length should be 16 + suffix
	expectedLen := 16 + len(suffix)
	if len(id) != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, len(id))
	}
}

func TestGenerateMessageID(t *testing.T) {
	id1 := GenerateMessageID()
	id2 := GenerateMessageID()

	// UUID should be 36 characters
	if len(id1) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(id1))
	}

	// Should contain hyphens at specific positions
	expectedHyphens := []int{8, 13, 18, 23}
	for _, pos := range expectedHyphens {
		if id1[pos] != '-' {
			t.Errorf("Expected hyphen at position %d, got %c", pos, id1[pos])
		}
	}

	// Should be unique
	if id1 == id2 {
		t.Error("Generated message IDs should be unique")
	}
}

func TestGenerateUniqueID(t *testing.T) {
	id1 := GenerateUniqueID()
	id2 := GenerateUniqueID()

	// Should be 8 characters
	if len(id1) != 8 {
		t.Errorf("Expected length 8, got %d", len(id1))
	}

	// Should be unique (very high probability)
	if id1 == id2 {
		t.Error("Generated unique IDs should be unique")
	}

	// Should contain only hexadecimal characters
	for _, char := range id1 {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			t.Errorf("ID should contain only hex characters, got %c", char)
		}
	}
}

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+628123456789", "628123456789"}, // Remove +
		{"08123456789", "628123456789"},   // Add 62 prefix
		{"8123456789", "628123456789"},    // Add 62 prefix (no leading 0)
		{"628123456789", "628123456789"},  // Already correct
		{"+658123456789", "658123456789"}, // Singapore, remove +
		{"658123456789", "658123456789"},  // Singapore, already correct
		{"+18123456789", "18123456789"},   // US, remove +
		{"18123456789", "18123456789"},    // US, already correct
		{"23456789", "6223456789"},        // No country code, default to 62
		{"", "62"},                        // Empty string gets 62
	}

	for _, test := range tests {
		result := NormalizePhoneNumber(test.input)
		if result != test.expected {
			t.Errorf("NormalizePhoneNumber(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}
