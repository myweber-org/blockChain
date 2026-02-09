
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func NormalizeUsername(username string) (string, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return "", errors.New("username cannot be empty")
	}
	if len(trimmed) < 3 {
		return "", errors.New("username must be at least 3 characters")
	}
	return strings.ToLower(trimmed), nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	normalizedUsername, err := NormalizeUsername(profile.Username)
	if err != nil {
		return UserProfile{}, err
	}
	profile.Username = normalizedUsername

	if err := ValidateEmail(profile.Email); err != nil {
		return UserProfile{}, err
	}

	if profile.Age < 0 || profile.Age > 150 {
		return UserProfile{}, errors.New("age must be between 0 and 150")
	}

	return profile, nil
}