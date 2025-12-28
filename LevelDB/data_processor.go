
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
	if len(rawData) == 0 {
		return ValidationError{Field: "data", Message: "empty input"}
	}

	if err := json.Unmarshal(rawData, target); err != nil {
		return ValidationError{Field: "structure", Message: err.Error()}
	}

	return nil
}

func ValidateStringField(value string, fieldName string, minLength int) error {
	if len(value) < minLength {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be at least %d characters", minLength),
		}
	}
	return nil
}

func ValidateNumericField(value float64, fieldName string, minValue float64) error {
	if value < minValue {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be greater than %.2f", minValue),
		}
	}
	return nil
}