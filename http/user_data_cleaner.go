
package main

import (
    "fmt"
    "regexp"
    "strings"
)

type UserData struct {
    Email       string
    PhoneNumber string
    Username    string
}

func SanitizeEmail(email string) (string, error) {
    email = strings.ToLower(strings.TrimSpace(email))
    pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
    matched, err := regexp.MatchString(pattern, email)
    if err != nil {
        return "", fmt.Errorf("email regex validation failed: %v", err)
    }
    if !matched {
        return "", fmt.Errorf("invalid email format")
    }
    return email, nil
}

func FormatPhoneNumber(phone string) string {
    re := regexp.MustCompile(`\D`)
    cleaned := re.ReplaceAllString(phone, "")
    if len(cleaned) == 10 {
        return fmt.Sprintf("(%s) %s-%s", cleaned[0:3], cleaned[3:6], cleaned[6:10])
    }
    return phone
}

func CleanUsername(username string) string {
    username = strings.TrimSpace(username)
    re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
    return re.ReplaceAllString(username, "")
}

func ProcessUserData(data UserData) (UserData, error) {
    var err error
    result := data

    result.Email, err = SanitizeEmail(data.Email)
    if err != nil {
        return UserData{}, fmt.Errorf("email sanitization failed: %v", err)
    }

    result.PhoneNumber = FormatPhoneNumber(data.PhoneNumber)
    result.Username = CleanUsername(data.Username)

    return result, nil
}

func main() {
    testData := UserData{
        Email:       "TEST@Example.COM",
        PhoneNumber: "123-456-7890",
        Username:    "  User@Name#123  ",
    }

    cleanedData, err := ProcessUserData(testData)
    if err != nil {
        fmt.Printf("Error processing user data: %v\n", err)
        return
    }

    fmt.Printf("Original: %+v\n", testData)
    fmt.Printf("Cleaned:  %+v\n", cleanedData)
}