
package main

import (
    "regexp"
    "strings"
)

type User struct {
    ID       int
    Username string
    Email    string
}

func ValidateUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
    trimmed := strings.TrimSpace(email)
    return strings.ToLower(trimmed)
}

func ValidateUserInput(user User) (bool, []string) {
    var errors []string

    if !ValidateUsername(user.Username) {
        errors = append(errors, "Username must be 3-20 characters and contain only letters, numbers, and underscores")
    }

    if user.Email == "" {
        errors = append(errors, "Email cannot be empty")
    }

    return len(errors) == 0, errors
}package main

import (
    "regexp"
    "strings"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

func ValidateUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
    trimmed := strings.TrimSpace(email)
    return strings.ToLower(trimmed)
}

func ValidateAge(age int) bool {
    return age >= 13 && age <= 120
}

func ProcessUserData(data UserData) (UserData, error) {
    if !ValidateUsername(data.Username) {
        return UserData{}, ErrInvalidUsername
    }

    sanitizedEmail := SanitizeEmail(data.Email)

    if !ValidateAge(data.Age) {
        return UserData{}, ErrInvalidAge
    }

    return UserData{
        Username: data.Username,
        Email:    sanitizedEmail,
        Age:      data.Age,
    }, nil
}

var (
    ErrInvalidUsername = errors.New("invalid username format")
    ErrInvalidAge      = errors.New("age must be between 13 and 120")
)