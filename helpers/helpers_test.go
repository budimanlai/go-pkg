package helpers

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Pointer Tests
// ============================================================================

func TestPointer(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		val := 42
		ptr := Pointer(val)
		if ptr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *ptr != val {
			t.Errorf("Expected %d, got %d", val, *ptr)
		}
	})

	t.Run("string", func(t *testing.T) {
		str := "hello"
		strPtr := Pointer(str)
		if strPtr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *strPtr != str {
			t.Errorf("Expected %s, got %s", str, *strPtr)
		}
	})

	t.Run("bool", func(t *testing.T) {
		val := true
		ptr := Pointer(val)
		if ptr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *ptr != val {
			t.Errorf("Expected %v, got %v", val, *ptr)
		}
	})

	t.Run("float", func(t *testing.T) {
		val := 3.14
		ptr := Pointer(val)
		if ptr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *ptr != val {
			t.Errorf("Expected %f, got %f", val, *ptr)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		person := Person{Name: "John", Age: 30}
		personPtr := Pointer(person)
		if personPtr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *personPtr != person {
			t.Errorf("Expected %+v, got %+v", person, *personPtr)
		}
	})

	t.Run("empty_string", func(t *testing.T) {
		str := ""
		ptr := Pointer(str)
		if ptr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *ptr != str {
			t.Errorf("Expected empty string, got %s", *ptr)
		}
	})

	t.Run("zero_value", func(t *testing.T) {
		val := 0
		ptr := Pointer(val)
		if ptr == nil {
			t.Fatal("Expected non-nil pointer")
		}
		if *ptr != 0 {
			t.Errorf("Expected 0, got %d", *ptr)
		}
	})
}

func TestDerefPointer(t *testing.T) {
	t.Run("non_nil_int", func(t *testing.T) {
		val := 42
		ptr := &val
		result := DerefPointer(ptr, 0)
		if result != val {
			t.Errorf("Expected %d, got %d", val, result)
		}
	})

	t.Run("nil_int", func(t *testing.T) {
		var nilPtr *int
		result := DerefPointer(nilPtr, 100)
		if result != 100 {
			t.Errorf("Expected default value 100, got %d", result)
		}
	})

	t.Run("non_nil_string", func(t *testing.T) {
		str := "world"
		strPtr := &str
		resultStr := DerefPointer(strPtr, "default")
		if resultStr != str {
			t.Errorf("Expected %s, got %s", str, resultStr)
		}
	})

	t.Run("nil_string", func(t *testing.T) {
		var nilStrPtr *string
		resultStr := DerefPointer(nilStrPtr, "fallback")
		if resultStr != "fallback" {
			t.Errorf("Expected fallback, got %s", resultStr)
		}
	})

	t.Run("non_nil_bool", func(t *testing.T) {
		val := true
		ptr := &val
		result := DerefPointer(ptr, false)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("nil_bool", func(t *testing.T) {
		var nilPtr *bool
		result := DerefPointer(nilPtr, false)
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("zero_value_vs_default", func(t *testing.T) {
		val := 0
		ptr := &val
		result := DerefPointer(ptr, 99)
		if result != 0 {
			t.Errorf("Expected 0 (actual value, not default), got %d", result)
		}
	})
}

// ============================================================================
// JSON Tests
// ============================================================================

func TestUnmarshalTo(t *testing.T) {
	t.Run("struct_success", func(t *testing.T) {
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
	})

	t.Run("map_success", func(t *testing.T) {
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
	})

	t.Run("array_success", func(t *testing.T) {
		jsonStr := `["one","two","three"]`
		result, err := UnmarshalTo[[]string](jsonStr)
		if err != nil {
			t.Fatalf("Failed to unmarshal array: %v", err)
		}

		if len(result) != 3 {
			t.Errorf("Expected array length 3, got %d", len(result))
		}
		if result[0] != "one" {
			t.Errorf("Expected first element 'one', got %s", result[0])
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
		}
		_, err := UnmarshalTo[Person]("invalid json")
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})

	t.Run("empty_json", func(t *testing.T) {
		result, err := UnmarshalTo[map[string]interface{}]("{}")
		if err != nil {
			t.Fatalf("Failed to unmarshal empty JSON: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty map, got %v", result)
		}
	})

	t.Run("nested_struct", func(t *testing.T) {
		jsonStr := `{"name":"Bob","address":{"city":"Jakarta","country":"Indonesia"}}`
		type Address struct {
			City    string `json:"city"`
			Country string `json:"country"`
		}
		type Person struct {
			Name    string  `json:"name"`
			Address Address `json:"address"`
		}

		person, err := UnmarshalTo[Person](jsonStr)
		if err != nil {
			t.Fatalf("Failed to unmarshal nested struct: %v", err)
		}

		if person.Address.City != "Jakarta" {
			t.Errorf("Expected city Jakarta, got %s", person.Address.City)
		}
	})
}

func TestUnmarshalFromMap(t *testing.T) {
	t.Run("struct_success", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		dataMap := map[string]interface{}{
			"name": "John",
			"age":  float64(30), // JSON numbers are float64
		}

		person, err := UnmarshalFromMap[Person](dataMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal from map: %v", err)
		}

		if person.Name != "John" {
			t.Errorf("Expected name John, got %s", person.Name)
		}
		if person.Age != 30 {
			t.Errorf("Expected age 30, got %d", person.Age)
		}
	})

	t.Run("empty_map", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
		}
		dataMap := map[string]interface{}{}

		person, err := UnmarshalFromMap[Person](dataMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal empty map: %v", err)
		}

		if person.Name != "" {
			t.Errorf("Expected empty name, got %s", person.Name)
		}
	})

	t.Run("nested_map", func(t *testing.T) {
		type Address struct {
			City string `json:"city"`
		}
		type Person struct {
			Name    string  `json:"name"`
			Address Address `json:"address"`
		}

		dataMap := map[string]interface{}{
			"name": "Alice",
			"address": map[string]interface{}{
				"city": "Bandung",
			},
		}

		person, err := UnmarshalFromMap[Person](dataMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal nested map: %v", err)
		}

		if person.Address.City != "Bandung" {
			t.Errorf("Expected city Bandung, got %s", person.Address.City)
		}
	})
}

// ============================================================================
// String Helper Tests
// ============================================================================

func TestGenerateTrxID(t *testing.T) {
	t.Run("length_and_format", func(t *testing.T) {
		id1 := GenerateTrxID()
		id2 := GenerateTrxID()

		// Check length (YYMMDDHHMMSS = 12 + 4 random = 16)
		if len(id1) != 16 {
			t.Errorf("Expected length 16, got %d", len(id1))
		}
		if len(id2) != 16 {
			t.Errorf("Expected length 16, got %d", len(id2))
		}

		// Check format (should be numeric)
		for _, char := range id1 {
			if char < '0' || char > '9' {
				t.Errorf("ID should contain only digits, got %c", char)
			}
		}
	})

	t.Run("uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := GenerateTrxID()
			if ids[id] {
				t.Errorf("Duplicate ID generated: %s", id)
			}
			ids[id] = true
			time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
		}
	})

	t.Run("timestamp_prefix", func(t *testing.T) {
		id := GenerateTrxID()
		now := time.Now()
		expectedPrefix := now.Format("060102")

		if !strings.HasPrefix(id, expectedPrefix) {
			t.Errorf("Expected ID to start with date %s, got %s", expectedPrefix, id[:6])
		}
	})
}

func TestGenerateTrxIDWithPrefix(t *testing.T) {
	t.Run("basic_prefix", func(t *testing.T) {
		prefix := "TXN"
		id := GenerateTrxIDWithPrefix(prefix)

		if !strings.HasPrefix(id, prefix) {
			t.Errorf("Expected ID to start with %s, got %s", prefix, id)
		}

		expectedLen := len(prefix) + 16
		if len(id) != expectedLen {
			t.Errorf("Expected length %d, got %d", expectedLen, len(id))
		}
	})

	t.Run("empty_prefix", func(t *testing.T) {
		id := GenerateTrxIDWithPrefix("")
		if len(id) != 16 {
			t.Errorf("Expected length 16 with empty prefix, got %d", len(id))
		}
	})

	t.Run("special_chars_prefix", func(t *testing.T) {
		prefix := "TRX-2024-"
		id := GenerateTrxIDWithPrefix(prefix)

		if !strings.HasPrefix(id, prefix) {
			t.Errorf("Expected ID to start with %s, got %s", prefix, id)
		}
	})
}

func TestGenerateTrxIDWithSuffix(t *testing.T) {
	t.Run("basic_suffix", func(t *testing.T) {
		suffix := "END"
		id := GenerateTrxIDWithSuffix(suffix)

		if !strings.HasSuffix(id, suffix) {
			t.Errorf("Expected ID to end with %s, got %s", suffix, id)
		}

		expectedLen := 16 + len(suffix)
		if len(id) != expectedLen {
			t.Errorf("Expected length %d, got %d", expectedLen, len(id))
		}
	})

	t.Run("empty_suffix", func(t *testing.T) {
		id := GenerateTrxIDWithSuffix("")
		if len(id) != 16 {
			t.Errorf("Expected length 16 with empty suffix, got %d", len(id))
		}
	})

	t.Run("special_chars_suffix", func(t *testing.T) {
		suffix := "-PROD"
		id := GenerateTrxIDWithSuffix(suffix)

		if !strings.HasSuffix(id, suffix) {
			t.Errorf("Expected ID to end with %s", suffix)
		}
	})
}

func TestGenerateMessageID(t *testing.T) {
	t.Run("uuid_format", func(t *testing.T) {
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
	})

	t.Run("multiple_generations", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10; i++ {
			id := GenerateMessageID()
			if ids[id] {
				t.Errorf("Duplicate message ID generated: %s", id)
			}
			ids[id] = true
		}
	})
}

func TestGenerateUniqueID(t *testing.T) {
	t.Run("length", func(t *testing.T) {
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
	})

	t.Run("hex_characters", func(t *testing.T) {
		id := GenerateUniqueID()

		// Should contain only hexadecimal characters
		for _, char := range id {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				t.Errorf("ID should contain only hex characters, got %c", char)
			}
		}
	})

	t.Run("uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := GenerateUniqueID()
			if ids[id] {
				t.Errorf("Duplicate unique ID generated: %s", id)
			}
			ids[id] = true
		}
	})
}

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"indonesian_with_plus", "+628123456789", "628123456789"},
		{"indonesian_with_zero", "08123456789", "628123456789"},
		{"indonesian_without_prefix", "8123456789", "628123456789"},
		{"indonesian_already_normalized", "628123456789", "628123456789"},
		{"singapore_with_plus", "+658123456789", "658123456789"},
		{"singapore_normalized", "658123456789", "658123456789"},
		{"us_with_plus", "+18123456789", "18123456789"},
		{"us_normalized", "18123456789", "18123456789"},
		{"short_number_indonesian", "23456789", "6223456789"},
		{"empty_string", "", "62"},
		{"only_zero", "0", "62"},
		{"indonesian_mobile_085", "085123456789", "6285123456789"},
		{"indonesian_mobile_081", "081234567890", "6281234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePhoneNumber(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizePhoneNumber_EdgeCases(t *testing.T) {
	t.Run("multiple_zeros", func(t *testing.T) {
		result := NormalizePhoneNumber("008123456789")
		if !strings.HasPrefix(result, "62") {
			t.Errorf("Expected result to start with 62, got %s", result)
		}
	})

	t.Run("canada_number", func(t *testing.T) {
		// Canada uses country code 1
		result := NormalizePhoneNumber("+14165551234")
		expected := "14165551234"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}
