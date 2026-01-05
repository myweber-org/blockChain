package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(username, "")
	if len(username) > 20 {
		username = username[:20]
	}
	return username
}

func TransformUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 120 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  user_123!@#","age":25}`
	processedData, err := TransformUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}package main

import (
	"errors"
	"strings"
)

func ValidateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if len(email) < 5 {
		return errors.New("email too short")
	}
	return nil
}

func TrimAndUpper(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToUpper(trimmed)
}

func FilterEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Bio      string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&;]`)
	return re.ReplaceAllString(input, "")
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ProcessUserData(data UserData) (UserData, error) {
	data.Username = SanitizeInput(data.Username)
	data.Email = SanitizeInput(data.Email)
	data.Bio = SanitizeInput(data.Bio)

	if !ValidateEmail(data.Email) {
		return data, &InvalidEmailError{Email: data.Email}
	}

	if len(data.Username) < 3 || len(data.Username) > 20 {
		return data, &InvalidUsernameError{Username: data.Username}
	}

	return data, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be between 3 and 20 characters: " + e.Username
}