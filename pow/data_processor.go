package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Username: transformedUsername,
		Email:    strings.TrimSpace(email),
		Age:      age,
	}
	err := ValidateUserData(userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

type UserData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func ValidateAndTransform(data []byte) (*UserData, error) {
    var user UserData
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, fmt.Errorf("invalid JSON format: %w", err)
    }

    user.Name = strings.TrimSpace(user.Name)
    if user.Name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }

    if !strings.Contains(user.Email, "@") {
        return nil, fmt.Errorf("invalid email format")
    }

    if user.Age < 0 || user.Age > 150 {
        return nil, fmt.Errorf("age must be between 0 and 150")
    }

    return &user, nil
}

func main() {
    rawData := []byte(`{"name":"John Doe","email":"john@example.com","age":30}`)
    processedUser, err := ValidateAndTransform(rawData)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Validated user: %+v\n", processedUser)
}