package main

import "fmt"

func calculateAverage(numbers []float64) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    var sum float64
    for _, num := range numbers {
        sum += num
    }
    
    return sum / float64(len(numbers))
}

func main() {
    data := []float64{10.5, 20.3, 30.7, 40.1, 50.9}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}package main

import (
	"errors"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	for _, r := range data.Username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	if !strings.Contains(data.Email, "@") || !strings.Contains(data.Email, ".") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ProcessUserInput(rawUsername, rawEmail string, age int) (UserData, error) {
	userData := UserData{
		Username: TransformUsername(rawUsername),
		Email:    NormalizeEmail(rawEmail),
		Age:      age,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || len(data) == 0 {
        return []float64{}
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := i; j < i+windowSize; j++ {
            sum += data[j]
        }
        result = append(result, sum/float64(windowSize))
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    averaged := movingAverage(sampleData, 3)
    fmt.Printf("Moving average result: %v\n", averaged)
}