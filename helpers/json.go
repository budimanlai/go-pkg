package helpers

import "encoding/json"

// UnmarshalTo deserializes a JSON string into a value of type T.
// It takes a JSON string as input and returns the unmarshaled value of type T
// along with any error that occurred during unmarshaling.
//
// Type parameter T can be any type that is compatible with json.Unmarshal.
//
// Parameters:
//   - jsonString: A string containing valid JSON data to be unmarshaled.
//
// Returns:
//   - T: The unmarshaled value of the specified type.
//   - error: An error if the JSON is invalid or cannot be unmarshaled into type T.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	person, err := UnmarshalTo[Person](`{"name":"John","age":30}`)
func UnmarshalTo[T any](jsonString string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

// UnmarshalFromMap deserializes a map[string]interface{} into a value of type T.
// It takes a map as input and returns the unmarshaled value of type T
// along with any error that occurred during unmarshaling.
//
// Type parameter T can be any type that is compatible with json.Unmarshal.
//
// Parameters:
//   - dataMap: A map[string]interface{} containing data to be unmarshaled.
//
// Returns:
//   - T: The unmarshaled value of the specified type.
//   - error: An error if the data cannot be marshaled/unmarshaled into type T.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	dataMap := map[string]interface{}{
//	    "name": "John",
//	    "age":  30,
//	}
//	person, err := UnmarshalFromMap[Person](dataMap)
func UnmarshalFromMap[T any](dataMap map[string]interface{}) (T, error) {
	var result T
	// First marshal the map to JSON bytes
	jsonBytes, err := json.Marshal(dataMap)
	if err != nil {
		return result, err
	}
	// Then unmarshal the JSON bytes into the target type
	err = json.Unmarshal(jsonBytes, &result)
	return result, err
}
