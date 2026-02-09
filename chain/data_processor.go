
package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type RecordValidator func([]string) error

func ProcessCSVData(input string, validator RecordValidator) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.TrimLeadingSpace = true

	var processedRecords [][]string
	lineNumber := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if validator != nil {
			if err := validator(record); err != nil {
				return nil, errors.New("validation failed at line " + string(rune(lineNumber)) + ": " + err.Error())
			}
		}

		processedRecords = append(processedRecords, record)
	}

	if len(processedRecords) == 0 {
		return nil, errors.New("no valid records found in CSV data")
	}

	return processedRecords, nil
}

func ValidateRecordLength(expected int) RecordValidator {
	return func(record []string) error {
		if len(record) != expected {
			return errors.New("record length mismatch")
		}
		return nil
	}
}