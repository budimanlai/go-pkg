package types

import (
	"fmt"
	"time"
)

// UTCTime is a custom time type that marshals to/from JSON in UTC RFC3339 format.
// When marshaled, it will always end with 'Z' to indicate UTC timezone.
type UTCTime time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t UTCTime) MarshalJSON() ([]byte, error) {
	// Konversi tipe kustom kembali ke time.Time
	regularTime := time.Time(t)

	// Format ke UTC dengan format RFC3339 (yang menghasilkan 'Z')
	formatted := regularTime.UTC().Format(time.RFC3339)

	// JSON string harus dalam tanda kutip, jadi kita tambahkan secara manual.
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

func (t UTCTime) String() string {
	// Format ke UTC dengan RFC3339, sama seperti di MarshalJSON
	return time.Time(t).UTC().Format(time.RFC3339)
}
