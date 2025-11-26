package helpers

import (
	"strings"
	"testing"
)

func TestNormalizePhoneNumber1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"indonesian_with_plus_and_hyphens", "+62-812-3456-789", "628123456789"},
		{"indonesian_with_zero_and_hyphens", "0812-3456-789", "628123456789"},
		{"indonesian_without_prefix_and_hyphens", "812-3456-789", "8123456789"},
		{"us_with_plus_and_hyphens", "+1-202-555-1234", "12025551234"},
		{"indonesian_with_plus", "+628123456789", "628123456789"},
		{"indonesian_with_zero", "08123456789", "628123456789"},
		{"indonesian_without_prefix", "8123456789", "8123456789"},
		{"indonesian_already_normalized", "628123456789", "628123456789"},
		{"singapore_with_plus", "+6581234567", "6581234567"},
		{"us_with_plus", "+12025551234", "12025551234"},
		{"indonesian_with_spaces", "0812 3456 789", "628123456789"},
		{"indonesian_with_parentheses", "(0812) 3456-789", "628123456789"},
		{"indonesian_with_dots", "0812.3456.789", "628123456789"},
		{"indonesian_mobile_085", "085123456789", "6285123456789"},
		{"indonesian_mobile_081", "081234567890", "6281234567890"},
		{"indonesian_mobile_089", "089876543210", "6289876543210"},
		{"empty_string", "", ""},
		{"only_plus", "+", ""},
		{"only_hyphens", "---", ""},
		{"mixed_special_chars", "+62(812)-345.6789", "628123456789"},
		{"indonesian_08_prefix", "08", "628"},
		{"indonesian_081", "081", "6281"},
		{"non_indonesian_starting_08", "0812", "62812"},
		{"letters_and_numbers", "abc081def2345ghi6789", "628123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePhoneNumber(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizePhoneNumber_EdgeCases1(t *testing.T) {
	t.Run("only_special_characters", func(t *testing.T) {
		result := NormalizePhoneNumber("+-().")
		if result != "" {
			t.Errorf("Expected empty string, got %s", result)
		}
	})

	t.Run("very_long_number", func(t *testing.T) {
		input := "08123456789012345678901234567890"
		result := NormalizePhoneNumber(input)
		if !strings.HasPrefix(result, "62") {
			t.Errorf("Expected result to start with 62, got %s", result)
		}
		expected := "628" + input[2:]
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("starts_with_multiple_zeros", func(t *testing.T) {
		result := NormalizePhoneNumber("008123456789")
		expected := "008123456789"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})

	t.Run("unicode_characters", func(t *testing.T) {
		result := NormalizePhoneNumber("０８１２３４５６７８９")
		if result != "" {
			t.Errorf("Expected empty string for unicode digits, got %s", result)
		}
	})

	t.Run("mixed_valid_invalid", func(t *testing.T) {
		result := NormalizePhoneNumber("0812!@#$3456%^&*789")
		expected := "628123456789"
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}
