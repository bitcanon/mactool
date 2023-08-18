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
package oui

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

type Oui struct {
	// The OUI assignment (for example "1A2B3C")
	Assignment string

	// The organization name
	Organization string

	// The organization street address
	Address string
}

// The OUI database
type OuiDb struct {
	// The OUI database
	Entries []Oui
}

func (db *OuiDb) FindOuiByAssignment(s string) *Oui {
	for _, entry := range db.Entries {
		if entry.Assignment == s {
			return &entry
		}
	}
	return nil
}

// LoadDatabase loads an OUI database in CSV format from the specified reader
func LoadDatabase(r io.Reader) (*OuiDb, error) {
	// Create a CSV reader
	reader := csv.NewReader(r)
	reader.Comma = ','

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Create the OUI database
	db := &OuiDb{}

	// Loop through the records
	for _, record := range records {
		// Create an OUI entry
		entry := Oui{
			Assignment:   record[1],
			Organization: record[2],
			Address:      record[3],
		}

		// Append the OUI entry to the database
		db.Entries = append(db.Entries, entry)
	}

	// Return the OUI database
	return db, nil
}

func DownloadDatabase(w io.Writer, url string) error {
	// Download the CSV database
	// url := "http://standards-oui.ieee.org/oui/oui.csv"

	// Perform the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check if the response status code indicates success
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download database file: %s", response.Status)
	}

	// Get the size of the file for progress calculation
	fileSize := response.ContentLength

	// Create a buffer for reading and copying data,
	// processing the data in chunks of 1024 bytes
	buf := make([]byte, 1024)

	var totalDownloaded int64

	// Read and write data in chunks, updating progress along the way
	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			// Write the data to the output file
			_, err := w.Write(buf[:n])
			if err != nil {
				return err
			}

			totalDownloaded += int64(n)

			// Calculate and display progress in percent
			progressPercent := (float64(totalDownloaded) / float64(fileSize)) * 100
			fmt.Printf("\rDownload Progress: %.2f%%", progressPercent)
		}

		if err != nil {
			// Check if we reached the end of the file
			if err == io.EOF {
				break
			}

			// Return an error if we encountered an error other than EOF
			return err
		}
	}

	// Print a newline after the progress indicator
	fmt.Println()

	return nil
}
