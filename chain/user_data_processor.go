package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

func ValidateUser(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON structure: %w", err)
	}

	if user.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if user.ID <= 0 {
		return nil, fmt.Errorf("ID must be positive integer")
	}

	return &user, nil
}

func ProcessUserJSON(input string) {
	user, err := ValidateUser([]byte(input))
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return
	}

	fmt.Printf("Valid user: %s (ID: %d, Email: %s, Active: %t)\n",
		user.Username, user.ID, user.Email, user.Active)
}

func main() {
	sampleJSON := `{
		"id": 42,
		"username": "johndoe",
		"email": "john@example.com",
		"active": true
	}`

	ProcessUserJSON(sampleJSON)
}