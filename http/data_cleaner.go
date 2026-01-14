package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
	Data string
}

func removeDuplicates(records []Record) []Record {
	seen := make(map[int]bool)
	var unique []Record

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func sortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
}

func main() {
	records := []Record{
		{ID: 3, Name: "Item C", Data: "Sample data C"},
		{ID: 1, Name: "Item A", Data: "Sample data A"},
		{ID: 2, Name: "Item B", Data: "Sample data B"},
		{ID: 1, Name: "Item A", Data: "Sample data A"},
		{ID: 4, Name: "Item D", Data: "Sample data D"},
	}

	fmt.Println("Original records:", len(records))
	uniqueRecords := removeDuplicates(records)
	sortRecords(uniqueRecords)

	fmt.Println("Unique records:", len(uniqueRecords))
	for _, r := range uniqueRecords {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}