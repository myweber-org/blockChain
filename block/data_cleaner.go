
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	duplicates map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		duplicates: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	unique := []string{}
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.duplicates[normalized] {
			dc.duplicates[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) Reset() {
	dc.duplicates = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"user@domain.org",
		"test@example.com",
		"invalid-email",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(emails)
	fmt.Println("Unique emails:", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		cleaned := make([]string, len(record))
		for i, field := range record {
			cleaned[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleaned); err != nil {
			return fmt.Errorf("error writing CSV: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := cleanCSV(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV cleaning completed successfully")
}package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read CSV record: %w", err)
        }

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("failed to write CSV record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    if err := cleanCSV(os.Args[1], os.Args[2]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV cleaning completed successfully")
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSVData(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	cleanedHeaders := make([]string, len(headers))
	for i, header := range headers {
		cleanedHeaders[i] = strings.TrimSpace(strings.ToLower(header))
	}

	if err := writer.Write(cleanedHeaders); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedField := strings.TrimSpace(field)
			if cleanedField == "" {
				cleanedField = "N/A"
			}
			cleanedRecord[i] = cleanedField
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	if err := cleanCSVData(inputPath, outputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaned successfully. Output saved to %s\n", outputPath)
}
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

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    headers, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("skipping line %d due to read error: %v\n", lineNum, err)
            continue
        }

        if len(row) != 4 {
            fmt.Printf("skipping line %d: invalid column count %d\n", lineNum, len(row))
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("skipping line %d: %v\n", lineNum, err)
            continue
        }

        if !validateRecord(record) {
            fmt.Printf("skipping line %d: validation failed\n", lineNum)
            continue
        }

        cleanedRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            fmt.Sprintf("%.2f", record.Score),
        }

        if err := writer.Write(cleanedRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var rec Record
    var err error

    rec.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return rec, fmt.Errorf("invalid ID: %w", err)
    }

    rec.Name = row[1]

    rec.Email = row[2]

    rec.Score, err = strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
    if err != nil {
        return rec, fmt.Errorf("invalid score: %w", err)
    }

    return rec, nil
}

func validateRecord(r Record) bool {
    if r.ID <= 0 {
        return false
    }
    if len(strings.TrimSpace(r.Name)) == 0 {
        return false
    }
    if !strings.Contains(r.Email, "@") {
        return false
    }
    if r.Score < 0 || r.Score > 100 {
        return false
    }
    return true
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}