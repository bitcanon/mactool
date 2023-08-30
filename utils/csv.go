package utils

import (
	"bytes"
	"encoding/csv"
)

// ConvertStringSliceToCSV converts a string slice to a CSV-formatted string
func ConvertStringSliceToCSV(data []string) (string, error) {
	// Create a buffer to write CSV data to
	var csvBuffer bytes.Buffer

	// If the data slice is empty, return an empty string
	if len(data) == 0 {
		return "", nil
	}

	// Create a new CSV writer that writes to the buffer
	csvWriter := csv.NewWriter(&csvBuffer)

	// Write the string slice as a CSV record using WriteAll
	err := csvWriter.WriteAll([][]string{data})
	if err != nil {
		return "", err
	}

	// Get the CSV-formatted string from the buffer
	csvString := csvBuffer.String()

	// Return the CSV-formatted string
	return csvString, nil
}
