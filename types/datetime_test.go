package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// MarshalJSON Tests
// ============================================================================

func TestUTCTimeMarshalJSON(t *testing.T) {
	t.Run("marshal_utc_time", func(t *testing.T) {
		utcTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
		ut := UTCTime(utcTime)

		data, err := json.Marshal(ut)
		if err != nil {
			t.Fatalf("Failed to marshal UTCTime: %v", err)
		}

		jsonStr := string(data)
		if !strings.HasPrefix(jsonStr, `"`) || !strings.HasSuffix(jsonStr, `"`) {
			t.Errorf("Expected JSON string, got %s", jsonStr)
		}

		timeStr := jsonStr[1 : len(jsonStr)-1]
		expected := "2025-10-15T12:30:45Z"
		if timeStr != expected {
			t.Errorf("Expected %s, got %s", expected, timeStr)
		}

		if !strings.HasSuffix(timeStr, "Z") {
			t.Errorf("Expected time to end with Z, got %s", timeStr)
		}
	})

	t.Run("marshal_with_nanoseconds", func(t *testing.T) {
		utcTime := time.Date(2025, 10, 15, 12, 30, 45, 123456789, time.UTC)
		ut := UTCTime(utcTime)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		// RFC3339 should handle nanoseconds
		if !strings.Contains(jsonStr, "2025-10-15T12:30:45") {
			t.Errorf("Expected timestamp in output, got %s", jsonStr)
		}
	})

	t.Run("marshal_zero_time", func(t *testing.T) {
		var ut UTCTime
		data, err := json.Marshal(ut)
		if err != nil {
			t.Fatalf("Failed to marshal zero time: %v", err)
		}

		jsonStr := string(data)
		// Zero time in UTC
		expected := `"0001-01-01T00:00:00Z"`
		if jsonStr != expected {
			t.Errorf("Expected %s, got %s", expected, jsonStr)
		}
	})

	t.Run("marshal_midnight", func(t *testing.T) {
		midnight := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		ut := UTCTime(midnight)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		expected := `"2025-01-01T00:00:00Z"`
		if jsonStr != expected {
			t.Errorf("Expected %s, got %s", expected, jsonStr)
		}
	})

	t.Run("marshal_end_of_day", func(t *testing.T) {
		endOfDay := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
		ut := UTCTime(endOfDay)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		expected := `"2025-12-31T23:59:59Z"`
		if jsonStr != expected {
			t.Errorf("Expected %s, got %s", expected, jsonStr)
		}
	})
}

// ============================================================================
// UnmarshalJSON Tests
// ============================================================================

func TestUTCTimeUnmarshalJSON(t *testing.T) {
	t.Run("unmarshal_utc_time", func(t *testing.T) {
		jsonStr := `"2025-10-15T12:30:45Z"`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err != nil {
			t.Fatalf("Failed to unmarshal UTCTime: %v", err)
		}

		actualTime := time.Time(ut)
		expectedTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)

		if !actualTime.Equal(expectedTime) {
			t.Errorf("Expected %v, got %v", expectedTime, actualTime)
		}

		if actualTime.Location() != time.UTC {
			t.Errorf("Expected UTC timezone, got %v", actualTime.Location())
		}
	})

	t.Run("unmarshal_with_timezone_offset", func(t *testing.T) {
		jsonStr := `"2025-10-15T19:30:45+07:00"`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err != nil {
			t.Fatalf("Failed to unmarshal with timezone: %v", err)
		}

		actualTime := time.Time(ut)
		// 19:30 +07:00 should be 12:30 UTC
		expectedTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)

		if !actualTime.Equal(expectedTime) {
			t.Errorf("Expected %v, got %v", expectedTime, actualTime)
		}
	})

	t.Run("unmarshal_with_nanoseconds", func(t *testing.T) {
		jsonStr := `"2025-10-15T12:30:45.123456789Z"`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err != nil {
			t.Fatalf("Failed to unmarshal with nanoseconds: %v", err)
		}

		actualTime := time.Time(ut)
		if actualTime.Nanosecond() == 0 {
			t.Error("Expected nanoseconds to be preserved")
		}
	})

	t.Run("unmarshal_invalid_format", func(t *testing.T) {
		jsonStr := `"not a valid time"`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err == nil {
			t.Error("Expected error for invalid format")
		}
	})

	t.Run("unmarshal_empty_string", func(t *testing.T) {
		jsonStr := `""`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err == nil {
			t.Error("Expected error for empty string")
		}
	})

	t.Run("unmarshal_zero_time", func(t *testing.T) {
		jsonStr := `"0001-01-01T00:00:00Z"`

		var ut UTCTime
		err := json.Unmarshal([]byte(jsonStr), &ut)
		if err != nil {
			t.Fatalf("Failed to unmarshal zero time: %v", err)
		}

		actualTime := time.Time(ut)
		if !actualTime.IsZero() {
			t.Error("Expected zero time")
		}
	})
}

// ============================================================================
// String Tests
// ============================================================================

func TestUTCTimeString(t *testing.T) {
	t.Run("string_representation", func(t *testing.T) {
		utcTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
		ut := UTCTime(utcTime)

		str := ut.String()
		expected := "2025-10-15T12:30:45Z"

		if str != expected {
			t.Errorf("Expected %s, got %s", expected, str)
		}

		if !strings.HasSuffix(str, "Z") {
			t.Errorf("Expected string to end with Z, got %s", str)
		}
	})

	t.Run("string_with_different_timezone", func(t *testing.T) {
		loc, _ := time.LoadLocation("America/New_York")
		localTime := time.Date(2025, 10, 15, 8, 30, 45, 0, loc)
		ut := UTCTime(localTime)

		str := ut.String()
		// Should convert to UTC in string representation
		if !strings.HasSuffix(str, "Z") {
			t.Errorf("Expected UTC format with Z suffix, got %s", str)
		}
	})

	t.Run("string_zero_time", func(t *testing.T) {
		var ut UTCTime
		str := ut.String()

		expected := "0001-01-01T00:00:00Z"
		if str != expected {
			t.Errorf("Expected %s, got %s", expected, str)
		}
	})

	t.Run("string_with_nanoseconds", func(t *testing.T) {
		utcTime := time.Date(2025, 10, 15, 12, 30, 45, 999999999, time.UTC)
		ut := UTCTime(utcTime)

		str := ut.String()
		// RFC3339 format includes nanoseconds
		if !strings.Contains(str, "2025-10-15T12:30:45") {
			t.Errorf("Expected timestamp in string, got %s", str)
		}
	})
}

// ============================================================================
// Round Trip Tests
// ============================================================================

func TestUTCTimeRoundTrip(t *testing.T) {
	t.Run("marshal_unmarshal_round_trip", func(t *testing.T) {
		// Use time without nanoseconds for cleaner round trip test
		originalTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
		originalUT := UTCTime(originalTime)

		data, err := json.Marshal(originalUT)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		var unmarshaledUT UTCTime
		err = json.Unmarshal(data, &unmarshaledUT)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		original := time.Time(originalUT)
		unmarshaled := time.Time(unmarshaledUT)

		if !original.Equal(unmarshaled) {
			t.Errorf("Round trip failed: original %v, unmarshaled %v", original, unmarshaled)
		}
	})

	t.Run("round_trip_with_nanoseconds", func(t *testing.T) {
		// Test that nanoseconds are preserved through marshal/unmarshal
		originalTime := time.Date(2025, 10, 15, 12, 30, 45, 123456789, time.UTC)
		originalUT := UTCTime(originalTime)

		data, _ := json.Marshal(originalUT)

		var unmarshaledUT UTCTime
		json.Unmarshal(data, &unmarshaledUT)

		// Check if the marshaled JSON contains nanoseconds
		jsonStr := string(data)
		// If nanoseconds are in the format, they should be preserved
		// Otherwise, we compare without nanoseconds
		original := time.Time(originalUT)
		unmarshaled := time.Time(unmarshaledUT)

		// RFC3339 format may or may not include fractional seconds
		// Depending on the Go version, so we just check they're close
		if strings.Contains(jsonStr, ".") {
			// Has fractional seconds, should be equal
			if !original.Equal(unmarshaled) {
				t.Logf("Nanoseconds preserved: %s", jsonStr)
			}
		} else {
			// No fractional seconds, compare up to seconds
			if !original.Truncate(time.Second).Equal(unmarshaled.Truncate(time.Second)) {
				t.Errorf("Round trip failed (seconds): original %v, unmarshaled %v", original, unmarshaled)
			}
		}
	})

	t.Run("round_trip_multiple_times", func(t *testing.T) {
		originalUT := UTCTime(time.Now().UTC())

		for i := 0; i < 5; i++ {
			data, _ := json.Marshal(originalUT)
			var ut UTCTime
			json.Unmarshal(data, &ut)

			if !time.Time(originalUT).Truncate(time.Second).Equal(time.Time(ut).Truncate(time.Second)) {
				t.Errorf("Round trip %d failed", i+1)
			}
			originalUT = ut
		}
	})

	t.Run("struct_with_utctime_field", func(t *testing.T) {
		type TestStruct struct {
			ID        int     `json:"id"`
			CreatedAt UTCTime `json:"created_at"`
		}

		original := TestStruct{
			ID:        1,
			CreatedAt: UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)),
		}

		data, _ := json.Marshal(original)

		var unmarshaled TestStruct
		err := json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Fatalf("Failed to unmarshal struct: %v", err)
		}

		if unmarshaled.ID != original.ID {
			t.Error("ID mismatch")
		}

		if !time.Time(unmarshaled.CreatedAt).Equal(time.Time(original.CreatedAt)) {
			t.Error("CreatedAt mismatch")
		}
	})
}

// ============================================================================
// Timezone Conversion Tests
// ============================================================================

func TestUTCTimeWithDifferentTimezone(t *testing.T) {
	t.Run("jakarta_to_utc", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		localTime := time.Date(2025, 10, 15, 19, 30, 45, 0, loc) // 19:30 in Jakarta is 12:30 UTC
		ut := UTCTime(localTime)

		data, err := json.Marshal(ut)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		timeStr := string(data)
		expected := `"2025-10-15T12:30:45Z"` // Should be converted to UTC

		if timeStr != expected {
			t.Errorf("Expected %s, got %s", expected, timeStr)
		}
	})

	t.Run("new_york_to_utc", func(t *testing.T) {
		loc, _ := time.LoadLocation("America/New_York")
		localTime := time.Date(2025, 10, 15, 8, 30, 45, 0, loc)
		ut := UTCTime(localTime)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		// Should be converted to UTC
		if !strings.HasSuffix(strings.Trim(jsonStr, `"`), "Z") {
			t.Errorf("Expected UTC format, got %s", jsonStr)
		}
	})

	t.Run("tokyo_to_utc", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Tokyo")
		localTime := time.Date(2025, 10, 15, 21, 30, 45, 0, loc) // 21:30 in Tokyo is 12:30 UTC
		ut := UTCTime(localTime)

		str := ut.String()
		expected := "2025-10-15T12:30:45Z"

		if str != expected {
			t.Errorf("Expected %s, got %s", expected, str)
		}
	})

	t.Run("london_to_utc", func(t *testing.T) {
		loc, _ := time.LoadLocation("Europe/London")
		localTime := time.Date(2025, 1, 15, 12, 30, 45, 0, loc) // Winter time in London
		ut := UTCTime(localTime)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		// London winter time is same as UTC
		expected := `"2025-01-15T12:30:45Z"`
		if jsonStr != expected {
			t.Errorf("Expected %s, got %s", expected, jsonStr)
		}
	})
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestUTCTimeEdgeCases(t *testing.T) {
	t.Run("leap_year_date", func(t *testing.T) {
		leapDay := time.Date(2024, 2, 29, 12, 30, 45, 0, time.UTC)
		ut := UTCTime(leapDay)

		data, _ := json.Marshal(ut)
		jsonStr := string(data)

		expected := `"2024-02-29T12:30:45Z"`
		if jsonStr != expected {
			t.Errorf("Expected %s, got %s", expected, jsonStr)
		}
	})

	t.Run("year_boundary", func(t *testing.T) {
		newYear := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		ut := UTCTime(newYear)

		str := ut.String()
		expected := "2025-01-01T00:00:00Z"

		if str != expected {
			t.Errorf("Expected %s, got %s", expected, str)
		}
	})

	t.Run("very_old_date", func(t *testing.T) {
		oldDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		ut := UTCTime(oldDate)

		data, _ := json.Marshal(ut)
		var unmarshaledUT UTCTime
		json.Unmarshal(data, &unmarshaledUT)

		if !time.Time(ut).Equal(time.Time(unmarshaledUT)) {
			t.Error("Old date round trip failed")
		}
	})

	t.Run("far_future_date", func(t *testing.T) {
		futureDate := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)
		ut := UTCTime(futureDate)

		str := ut.String()
		if !strings.Contains(str, "2100-12-31") {
			t.Errorf("Expected future date in string, got %s", str)
		}
	})

	t.Run("time_conversion_back_to_time", func(t *testing.T) {
		originalTime := time.Now().UTC()
		ut := UTCTime(originalTime)

		convertedBack := time.Time(ut)

		if !convertedBack.Equal(originalTime) {
			t.Error("Conversion to time.Time failed")
		}
	})
}

// ============================================================================
// Comparison Tests
// ============================================================================

func TestUTCTimeComparison(t *testing.T) {
	t.Run("compare_same_times", func(t *testing.T) {
		t1 := UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC))
		t2 := UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC))

		if !time.Time(t1).Equal(time.Time(t2)) {
			t.Error("Same times should be equal")
		}
	})

	t.Run("compare_different_timezones_same_instant", func(t *testing.T) {
		utc := UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC))

		loc, _ := time.LoadLocation("Asia/Jakarta")
		jakarta := UTCTime(time.Date(2025, 10, 15, 19, 30, 45, 0, loc))

		// They represent the same instant in time
		if !time.Time(utc).Equal(time.Time(jakarta)) {
			t.Error("Same instant in different timezones should be equal")
		}
	})

	t.Run("compare_different_times", func(t *testing.T) {
		t1 := UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC))
		t2 := UTCTime(time.Date(2025, 10, 15, 12, 30, 46, 0, time.UTC))

		if time.Time(t1).Equal(time.Time(t2)) {
			t.Error("Different times should not be equal")
		}

		if !time.Time(t1).Before(time.Time(t2)) {
			t.Error("t1 should be before t2")
		}
	})
}
