package main

import (
	"fmt"
	"sort"
)

type Record struct {
	ID   int
	Name string
}

func cleanData(records []Record) []Record {
	seen := make(map[int]bool)
	var unique []Record

	for _, r := range records {
		if !seen[r.ID] {
			seen[r.ID] = true
			unique = append(unique, r)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ID < unique[j].ID
	})

	return unique
}

func main() {
	data := []Record{
		{3, "Charlie"},
		{1, "Alice"},
		{2, "Bob"},
		{1, "Alice"},
		{4, "David"},
	}

	cleaned := cleanData(data)
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s\n", r.ID, r.Name)
	}
}