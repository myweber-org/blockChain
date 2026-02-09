
package main

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ValidateUserData(data UserData) error {
	if err := ValidateEmail(data.Email); err != nil {
		return err
	}

	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	normalizedUsername := NormalizeUsername(username)
	userData := UserData{
		Email:    email,
		Username: normalizedUsername,
		Age:      age,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}
package main

import (
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	return strings.TrimSpace(input)
}

func NormalizeWhitespace(input string) string {
	var result strings.Builder
	prevSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}

func RemoveSpecialChars(input string, keepSet string) string {
	var result strings.Builder
	keepMap := make(map[rune]bool)

	for _, r := range keepSet {
		keepMap[r] = true
	}

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || keepMap[r] {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func ProcessData(input string) string {
	cleaned := CleanInput(input)
	normalized := NormalizeWhitespace(cleaned)
	return RemoveSpecialChars(normalized, "._-@")
}