package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Tags      []string `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d: %s", user.ID, user.Email)
		}
		
		transformedUser := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     strings.ToLower(user.Email),
			"age_group": categorizeAge(user.Age),
			"status":    user.Active,
			"tag_count": len(user.Tags),
		}
		transformed = append(transformed, transformedUser)
	}
	
	return transformed, nil
}

func categorizeAge(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age < 30:
		return "young_adult"
	case age >= 30 && age < 50:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserJSON(jsonData []byte) ([]map[string]interface{}, error) {
	var users []UserProfile
	
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	
	activeUsers := FilterInactiveUsers(users)
	return TransformUserData(activeUsers)
}

func main() {
	jsonData := []byte(`[
		{"id":1,"username":"JohnDoe","email":"john@example.com","age":25,"active":true,"tags":["golang","backend"]},
		{"id":2,"username":"JaneSmith","email":"jane@test.org","age":32,"active":false,"tags":["frontend"]},
		{"id":3,"username":"BobWilson","email":"bob@domain.co","age":17,"active":true,"tags":[]}
	]`)
	
	result, err := ProcessUserJSON(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
}