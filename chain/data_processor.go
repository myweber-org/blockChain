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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func NormalizeUserData(data UserData) UserData {
	return UserData{
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Username: strings.TrimSpace(data.Username),
		Age:      data.Age,
	}
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	data := UserData{
		Email:    email,
		Username: username,
		Age:      age,
	}
	
	normalizedData := NormalizeUserData(data)
	
	if err := ValidateUserData(normalizedData); err != nil {
		return UserData{}, err
	}
	
	return normalizedData, nil
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email     string
	Username  string
	Age       int
	Biography string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(profile.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if profile.Age < 0 || profile.Age > 120 {
		return errors.New("age must be between 0 and 120")
	}

	if len(profile.Biography) > 500 {
		return errors.New("biography cannot exceed 500 characters")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(strings.TrimSpace(profile.Username))
	transformed.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	transformed.Biography = strings.TrimSpace(profile.Biography)

	if transformed.Biography == "" {
		transformed.Biography = "No biography provided."
	}

	return transformed
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, err
	}

	return TransformProfile(profile), nil
}