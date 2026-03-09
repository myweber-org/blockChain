package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func FilterCSVRows(input io.Reader, output io.Writer, filterFunc func([]string) bool) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filterFunc(record) {
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	return nil
}

func DefaultFilter(record []string) bool {
	for _, field := range record {
		trimmed := strings.TrimSpace(field)
		if trimmed == "" || strings.Contains(trimmed, "NULL") {
			return false
		}
	}
	return true
}