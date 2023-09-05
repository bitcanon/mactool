/*
Copyright © 2023 Mikael Schultz <bitcanon@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
