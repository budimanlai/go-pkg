package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUTCTimeMarshalJSON(t *testing.T) {
	// Create a specific time in UTC
	utcTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
	ut := UTCTime(utcTime)

	// Marshal to JSON
	data, err := json.Marshal(ut)
	if err != nil {
		t.Fatalf("Failed to marshal UTCTime: %v", err)
	}

	// Check that it's a JSON string
	jsonStr := string(data)
	if !strings.HasPrefix(jsonStr, `"`) || !strings.HasSuffix(jsonStr, `"`) {
		t.Errorf("Expected JSON string, got %s", jsonStr)
	}

	// Remove quotes and check format
	timeStr := jsonStr[1 : len(jsonStr)-1]
	expected := "2025-10-15T12:30:45Z"
	if timeStr != expected {
		t.Errorf("Expected %s, got %s", expected, timeStr)
	}

	// Ensure it ends with Z (UTC)
	if !strings.HasSuffix(timeStr, "Z") {
		t.Errorf("Expected time to end with Z, got %s", timeStr)
	}
}

func TestUTCTimeUnmarshalJSON(t *testing.T) {
	// JSON string representing UTC time
	jsonStr := `"2025-10-15T12:30:45Z"`

	var ut UTCTime
	err := json.Unmarshal([]byte(jsonStr), &ut)
	if err != nil {
		t.Fatalf("Failed to unmarshal UTCTime: %v", err)
	}

	// Convert back to time.Time and check
	actualTime := time.Time(ut)
	expectedTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)

	if !actualTime.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, actualTime)
	}

	// Ensure it's in UTC
	if actualTime.Location() != time.UTC {
		t.Errorf("Expected UTC timezone, got %v", actualTime.Location())
	}
}

func TestUTCTimeString(t *testing.T) {
	// Create a specific time
	utcTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
	ut := UTCTime(utcTime)

	// Get string representation
	str := ut.String()
	expected := "2025-10-15T12:30:45Z"

	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}

	// Ensure it ends with Z
	if !strings.HasSuffix(str, "Z") {
		t.Errorf("Expected string to end with Z, got %s", str)
	}
}

func TestUTCTimeRoundTrip(t *testing.T) {
	// Test marshal -> unmarshal round trip
	originalTime := time.Date(2025, 10, 15, 12, 30, 45, 123456789, time.UTC)
	originalUT := UTCTime(originalTime)

	// Marshal
	data, err := json.Marshal(originalUT)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaledUT UTCTime
	err = json.Unmarshal(data, &unmarshaledUT)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare (note: nanoseconds might be truncated in RFC3339 format)
	original := time.Time(originalUT)
	unmarshaled := time.Time(unmarshaledUT)

	// RFC3339 doesn't include nanoseconds, so compare up to seconds
	if !original.Truncate(time.Second).Equal(unmarshaled.Truncate(time.Second)) {
		t.Errorf("Round trip failed: original %v, unmarshaled %v", original, unmarshaled)
	}
}

func TestUTCTimeWithDifferentTimezone(t *testing.T) {
	// Create time in a different timezone
	loc, _ := time.LoadLocation("Asia/Jakarta")
	localTime := time.Date(2025, 10, 15, 19, 30, 45, 0, loc) // 19:30 in Jakarta is 12:30 UTC
	ut := UTCTime(localTime)

	// Marshal should convert to UTC
	data, err := json.Marshal(ut)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	timeStr := string(data)
	expected := `"2025-10-15T12:30:45Z"` // Should be converted to UTC

	if timeStr != expected {
		t.Errorf("Expected %s, got %s", expected, timeStr)
	}
}
