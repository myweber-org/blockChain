
package main

import (
	"regexp"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Bio      string
}

func SanitizeUsername(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]")
	sanitized := reg.ReplaceAllString(input, "")
	return strings.ToLower(sanitized)
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func TrimAndCleanBio(bio string, maxLength int) string {
	bio = strings.TrimSpace(bio)
	var cleaned strings.Builder
	for _, r := range bio {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			cleaned.WriteRune(r)
		}
	}
	result := cleaned.String()
	if len(result) > maxLength {
		result = result[:maxLength]
	}
	return result
}

func CleanUserData(user *UserData) error {
	user.Username = SanitizeUsername(user.Username)
	if !ValidateEmail(user.Email) {
		return ErrInvalidEmail
	}
	user.Bio = TrimAndCleanBio(user.Bio, 500)
	return nil
}

var ErrInvalidEmail = errors.New("invalid email format")