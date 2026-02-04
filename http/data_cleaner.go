
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func deduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    result := []DataRecord{}
    
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            result = append(result, record)
        }
    }
    return result
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanPhoneNumber(phone string) string {
    cleaned := strings.Builder{}
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            cleaned.WriteRune(ch)
        }
    }
    return cleaned.String()
}

func processRecords(records []DataRecord) []DataRecord {
    deduped := deduplicateRecords(records)
    
    for i := range deduped {
        deduped[i].Phone = cleanPhoneNumber(deduped[i].Phone)
    }
    
    validRecords := []DataRecord{}
    for _, record := range deduped {
        if validateEmail(record.Email) && len(record.Phone) >= 10 {
            validRecords = append(validRecords, record)
        }
    }
    
    return validRecords
}

func main() {
    sampleData := []DataRecord{
        {ID: 1, Email: "test@example.com", Phone: "(123) 456-7890"},
        {ID: 2, Email: "invalid-email", Phone: "555-1234"},
        {ID: 1, Email: "test@example.com", Phone: "1234567890"},
        {ID: 3, Email: "user@domain.org", Phone: "+1-800-555-0199"},
    }
    
    cleaned := processRecords(sampleData)
    
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}