
package data_processor

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, target interface{}) error {
	if err := json.Unmarshal(rawData, target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if validator, ok := target.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func (u *UserData) Validate() error {
	if u.Username == "" {
		return ValidationError{Field: "username", Message: "cannot be empty"}
	}
	if u.Email == "" {
		return ValidationError{Field: "email", Message: "cannot be empty"}
	}
	if u.Age < 0 || u.Age > 150 {
		return ValidationError{Field: "age", Message: "must be between 0 and 150"}
	}
	return nil
}