package oui

import (
	"encoding/csv"
	"io"
)

type Oui struct {
	// The OUI assignment
	Assignment string
	// The organization name
	Organization string
	// The organization address
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
