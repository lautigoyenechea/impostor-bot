package main

import (
	"encoding/csv"
	"os"
)

type ImpostorWord struct {
	Word    string
	Concept string
	Hint    string
}

func LoadImpostorWords(path string) ([]ImpostorWord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var words []ImpostorWord

	// Skip header
	for i := 1; i < len(records); i++ {
		record := records[i]

		words = append(words, ImpostorWord{
			Word:    record[0],
			Concept: record[1],
			Hint:    record[2],
		})
	}

	return words, nil
}
