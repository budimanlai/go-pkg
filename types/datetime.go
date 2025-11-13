package types

import (
	"fmt"
	"time"
)

// UTCTime is a custom time type that provides automatic UTC conversion for JSON serialization.
// It ensures all time values are stored and transmitted in UTC timezone with RFC3339 format.
//
// Features:
//   - Always marshals to UTC with 'Z' suffix (e.g., "2025-10-15T04:56:56Z")
//   - Automatically converts from any timezone to UTC during marshaling
//   - Implements json.Marshaler and json.Unmarshaler interfaces
//   - Provides String() method for readable output
//
// Use cases:
//   - API responses requiring consistent UTC timestamps
//   - Database models where timezone normalization is needed
//   - Any scenario requiring timezone-agnostic time handling
//
// Example:
//
//	type User struct {
//	    ID        int     `json:"id"`
//	    CreatedAt UTCTime `json:"created_at"`
//	    UpdatedAt UTCTime `json:"updated_at"`
//	}
//
//	user := User{
//	    ID:        1,
//	    CreatedAt: UTCTime(time.Now()),
//	    UpdatedAt: UTCTime(time.Now()),
//	}
//	// When marshaled to JSON:
//	// {"id":1,"created_at":"2025-10-15T04:56:56Z","updated_at":"2025-10-15T04:56:56Z"}
type UTCTime time.Time

// MarshalJSON implements the json.Marshaler interface for UTCTime.
// It converts the time to UTC and formats it as RFC3339 with 'Z' suffix.
//
// The output format is always: "YYYY-MM-DDTHH:MM:SSZ"
// Regardless of the original timezone, the time is converted to UTC before marshaling.
//
// Returns:
//   - []byte: JSON-encoded string with quotes (e.g., []byte(`"2025-10-15T04:56:56Z"`))
//   - error: Error if formatting fails (rarely occurs)
//
// Example:
//
//	t := UTCTime(time.Date(2025, 10, 15, 4, 56, 56, 0, time.UTC))
//	json, _ := t.MarshalJSON()
//	// Output: []byte(`"2025-10-15T04:56:56Z"`)
func (t UTCTime) MarshalJSON() ([]byte, error) {
	// Konversi tipe kustom kembali ke time.Time
	regularTime := time.Time(t)

	// Format ke UTC dengan format RFC3339 (yang menghasilkan 'Z')
	formatted := regularTime.UTC().Format(time.RFC3339)

	// JSON string harus dalam tanda kutip, jadi kita tambahkan secara manual.
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for UTCTime.
// It parses a JSON time string in RFC3339 format and converts it to UTCTime.
//
// Accepted formats:
//   - "2025-10-15T04:56:56Z" (UTC with Z suffix)
//   - "2025-10-15T04:56:56+07:00" (with timezone offset)
//   - Any valid RFC3339 format
//
// Parameters:
//   - data: JSON-encoded time string with quotes (e.g., []byte(`"2025-10-15T04:56:56Z"`))
//
// Returns:
//   - error: Error if the data cannot be parsed as RFC3339 format
//
// Example:
//
//	var t UTCTime
//	data := []byte(`"2025-10-15T04:56:56Z"`)
//	err := t.UnmarshalJSON(data)
//	// t now contains the parsed time
func (t *UTCTime) UnmarshalJSON(data []byte) error {
	// data adalah string JSON dengan tanda kutip, misal: []byte(`"2025-10-15T04:56:56Z"`)
	// Kita perlu menghapus tanda kutip sebelum mem-parsing.
	// time.RFC3339 sudah mengharapkan format seperti itu.
	parsedTime, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err != nil {
		return err
	}
	*t = UTCTime(parsedTime)
	return nil
}

// String returns a string representation of UTCTime in RFC3339 format with UTC timezone.
// This method is useful for logging, debugging, and displaying time values.
//
// The output format is always: "YYYY-MM-DDTHH:MM:SSZ"
//
// Returns:
//   - string: Time formatted as RFC3339 in UTC (e.g., "2025-10-15T04:56:56Z")
//
// Example:
//
//	t := UTCTime(time.Now())
//	fmt.Println(t.String())
//	// Output: 2025-10-15T04:56:56Z
func (t UTCTime) String() string {
	// Format ke UTC dengan RFC3339, sama seperti di MarshalJSON
	return time.Time(t).UTC().Format(time.RFC3339)
}
