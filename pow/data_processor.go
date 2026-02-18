
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	return strings.Title(name)
}

func ProcessUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Name = SanitizeName(data.Name)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","name":" john doe ","age":30}`
	processed, err := ProcessUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}