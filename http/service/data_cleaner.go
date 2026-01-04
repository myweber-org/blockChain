
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID    int
    Name  string
    Email string
    Score float64
}

func cleanEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}

func validateScore(score float64) bool {
    return score >= 0 && score <= 100
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNum := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNum++
        if lineNum == 1 {
            continue
        }

        if len(line) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(line[0]))
        if err != nil {
            continue
        }

        name := strings.TrimSpace(line[1])
        email := cleanEmail(line[2])

        score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
        if err != nil || !validateScore(score) {
            continue
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Email: email,
            Score: score,
        })
    }

    return records, nil
}

func calculateAverageScore(records []Record) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0.0
    for _, record := range records {
        total += record.Score
    }
    return total / float64(len(records))
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_cleaner <csv_file>")
        return
    }

    records, err := parseCSVFile(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    fmt.Printf("Processed %d valid records\n", len(records))
    fmt.Printf("Average score: %.2f\n", calculateAverageScore(records))

    for i, record := range records {
        if i < 3 {
            fmt.Printf("Sample record: %+v\n", record)
        }
    }
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := normalizeSpaces(trimmed)
	cleaned := removeSpecialChars(normalized)
	return cleaned
}

func normalizeSpaces(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}

func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}