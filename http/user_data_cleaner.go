package main

import (
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email       string
	PhoneNumber string
}

func CleanUserData(data UserData) (UserData, error) {
	cleaned := UserData{}

	if err := validateEmail(data.Email); err != nil {
		return cleaned, fmt.Errorf("email validation failed: %w", err)
	}
	cleaned.Email = strings.ToLower(strings.TrimSpace(data.Email))

	formattedPhone, err := formatPhoneNumber(data.PhoneNumber)
	if err != nil {
		return cleaned, fmt.Errorf("phone formatting failed: %w", err)
	}
	cleaned.PhoneNumber = formattedPhone

	return cleaned, nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func formatPhoneNumber(phone string) (string, error) {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	if len(digits) != 10 {
		return "", fmt.Errorf("phone number must contain 10 digits")
	}

	return fmt.Sprintf("(%s) %s-%s", digits[:3], digits[3:6], digits[6:]), nil
}

func main() {
	user := UserData{
		Email:       "TEST@Example.COM",
		PhoneNumber: "123-456-7890",
	}

	cleaned, err := CleanUserData(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Cleaned Email: %s\n", cleaned.Email)
	fmt.Printf("Formatted Phone: %s\n", cleaned.PhoneNumber)
}