package helpers

import "time"

// StringToDate converts a string to a date format (YYYY-MM-DD).
//
// Parameters:
//   - dateStr: String representing the date in "YYYY-MM-DD" format
//
// Returns:
//   - string: Date in "YYYY-MM-DD" format
//   - error: Error if parsing fails
//
// Example:
//
//	date, err := StringToDate("2024-11-13")
//	// Output: "2024-11-13", nil
func StringToDate(dateStr string) (time.Time, error) {
	const layout = "2006-01-02"
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// FormatDatetime formats a time.Time to "YYYY-MM-DD HH:MM:SS" using the time's
// own timezone (no UTC conversion). Returns empty string for zero time.
//
// Example:
//
//	t, _ := time.Parse(time.RFC3339, "2026-06-13T12:40:54.840622+07:00")
//	FormatDatetime(t) // "2026-06-13 12:40:54"
func FormatDatetime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDatetimePtr is a nil-safe variant of FormatDatetime.
// Returns empty string when t is nil or zero.
func FormatDatetimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatDatetime(*t)
}
