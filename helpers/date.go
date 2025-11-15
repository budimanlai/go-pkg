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
