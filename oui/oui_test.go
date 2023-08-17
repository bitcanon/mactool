package oui_test

import (
	"strings"
	"testing"

	"github.com/bitcanon/mactool/oui"
)

func TestLoadDatabase(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,583653,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,58A15F,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Verify that the database was loaded correctly
	if len(db.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(db.Entries))
	}

	// Verify that the first entry was loaded correctly
	if db.Entries[0].Assignment != "583653" {
		t.Errorf("expected 583653, got %s", db.Entries[0].Assignment)
	}

	// Verify that the Organization name of the first entry was loaded correctly
	if db.Entries[0].Organization != "Apple, Inc." {
		t.Errorf("expected Apple, Inc., got %s", db.Entries[0].Organization)
	}

	// Verify that the Organization address of the first entry was loaded correctly
	if db.Entries[0].Address != "1 Infinite Loop Cupertino CA US 95014" {
		t.Errorf("expected 1 Infinite Loop Cupertino CA US 95014, got %s", db.Entries[0].Address)
	}

	// Verify that the second entry was loaded correctly
	if db.Entries[1].Assignment != "58A15F" {
		t.Errorf("expected 58A15F, got %s", db.Entries[1].Assignment)
	}
}

func TestFindOui(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,583653,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,58A15F,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Find the OUI assignment
	oui := db.FindOuiByAssignment("583653")
	if oui == nil {
		t.Errorf("expected oui, got nil")
	}

	// Verify that the OUI assignment was found
	if oui.Assignment != "583653" {
		t.Errorf("expected 583653, got %s", oui.Assignment)
	}

	// Verify that the Organization name of the OUI assignment was found
	if oui.Organization != "Apple, Inc." {
		t.Errorf("expected Apple, Inc., got %s", oui.Organization)
	}

	// Verify that the Organization address of the OUI assignment was found
	if oui.Address != "1 Infinite Loop Cupertino CA US 95014" {
		t.Errorf("expected 1 Infinite Loop Cupertino CA US 95014, got %s", oui.Address)
	}
}

func TestFindOuiNotFound(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,583653,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,58A15F,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Find the OUI assignment
	oui := db.FindOuiByAssignment("000000")
	if oui != nil {
		t.Errorf("expected nil, got oui")
	}
}
