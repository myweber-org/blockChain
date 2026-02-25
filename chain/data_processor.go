
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// DataPayload represents a simple JSON structure
type DataPayload struct {
    ID    int    `json:"id"`
    Value string `json:"value"`
    Valid bool   `json:"valid"`
}

// ParseAndValidateJSON parses a JSON string into DataPayload and performs basic validation
func ParseAndValidateJSON(input string) (*DataPayload, error) {
    if strings.TrimSpace(input) == "" {
        return nil, fmt.Errorf("input string is empty")
    }

    var payload DataPayload
    err := json.Unmarshal([]byte(input), &payload)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }

    if payload.ID <= 0 {
        return nil, fmt.Errorf("invalid ID: must be positive integer")
    }

    if payload.Value == "" {
        return nil, fmt.Errorf("value field cannot be empty")
    }

    return &payload, nil
}

func main() {
    testJSON := `{"id": 123, "value": "test data", "valid": true}`
    result, err := ParseAndValidateJSON(testJSON)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Parsed payload: %+v\n", result)
}