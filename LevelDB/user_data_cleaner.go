package main

import (
    "fmt"
    "regexp"
    "strings"
)

type UserData struct {
    Email    string
    Phone    string
    Username string
}

func cleanEmail(email string) string {
    trimmed := strings.TrimSpace(email)
    return strings.ToLower(trimmed)
}

func formatPhoneNumber(phone string) string {
    re := regexp.MustCompile(`\D`)
    cleaned := re.ReplaceAllString(phone, "")
    
    if len(cleaned) == 10 {
        return fmt.Sprintf("(%s) %s-%s", 
            cleaned[0:3], 
            cleaned[3:6], 
            cleaned[6:10])
    }
    return cleaned
}

func validateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

func CleanUserData(data UserData) UserData {
    cleaned := UserData{
        Email:    cleanEmail(data.Email),
        Phone:    formatPhoneNumber(data.Phone),
        Username: strings.TrimSpace(data.Username),
    }
    
    if !validateEmail(cleaned.Email) {
        cleaned.Email = ""
    }
    
    return cleaned
}

func main() {
    sampleData := UserData{
        Email:    "  TEST@Example.COM  ",
        Phone:    "123-456-7890",
        Username: "  john_doe  ",
    }
    
    cleaned := CleanUserData(sampleData)
    
    fmt.Printf("Original: %+v\n", sampleData)
    fmt.Printf("Cleaned:  %+v\n", cleaned)
}