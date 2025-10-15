package main

import (
	"encoding/json"
	"fmt"
	"time"

	types "github.com/budimanlai/go-pkg/types"
)

type Users struct {
	ID        int           `json:"id"`
	CreatedAt types.UTCTime `json:"created_at"`
}

var (
	DateTimeFormat = time.RFC3339
)

func main() {
	user := Users{
		ID:        1,
		CreatedAt: types.UTCTime(time.Now()),
	}

	fmt.Println("Formatted CreatedAt:", user.CreatedAt)

	// parsing user ke json string
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON Data:", string(jsonData))

	jsonString := `{"id":100,"created_at":"2025-10-01T15:56:56Z"}`
	var user2 Users
	err = json.Unmarshal([]byte(jsonString), &user2)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("")
	fmt.Println("JSON string:", jsonString)
	fmt.Println("Parsed User:", user2)
	fmt.Println("Parsed CreatedAt:", user2.CreatedAt)
}
