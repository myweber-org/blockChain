package data_processor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidateJSONString checks if the provided string is valid JSON.
func ValidateJSONString(input string) (bool, error) {
	var js interface{}
	err := json.Unmarshal([]byte(input), &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ExtractStringField safely extracts a string field from a map.
func ExtractStringField(data map[string]interface{}, key string) (string, error) {
	val, exists := data[key]
	if !exists {
		return "", fmt.Errorf("key '%s' not found", key)
	}

	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value for key '%s' is not a string", key)
	}

	return strings.TrimSpace(strVal), nil
}

// FlattenMapToString flattens a map into a single string for logging.
func FlattenMapToString(m map[string]string) string {
	var builder strings.Builder
	first := true
	for k, v := range m {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		first = false
	}
	return builder.String()
}