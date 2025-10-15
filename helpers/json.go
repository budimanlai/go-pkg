package helpers

import "encoding/json"

// UnmarshalTo unmarshals a JSON string into a struct of type T.
func UnmarshalTo[T any](jsonString string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	return result, err
}
