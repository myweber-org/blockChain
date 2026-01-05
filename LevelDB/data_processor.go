
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func ValidateAge(age int) error {
	if age < 0 || age > 120 {
		return errors.New("age must be between 0 and 120")
	}
	return nil
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}

	data.Username = SanitizeUsername(data.Username)

	if err := ValidateAge(data.Age); err != nil {
		return UserData{}, err
	}

	return data, nil
}
package main

import (
	"errors"
	"strings"
)

func ValidateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if strings.Count(email, "@") > 1 {
		return errors.New("invalid email format")
	}
	return nil
}

func TrimAndTitle(s string) string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return trimmed
	}
	return strings.ToUpper(trimmed[:1]) + strings.ToLower(trimmed[1:])
}

func FilterEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}