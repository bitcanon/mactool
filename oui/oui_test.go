package oui_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitcanon/mactool/oui"
)

// TestLoadDatabase tests the case where the database is loaded successfully.
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

// TestFindOui tests the case where the OUI assignment is found.
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

// TestFindOuiNotFound tests the case where the OUI assignment is not found.
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

// TestDownloadDatabase tests the case where the HTTP server returns a valid
// database. DownloadDatabase() should return the database in this case.
func TestDownloadDatabase(t *testing.T) {
	// Create a mock HTTP server for successful downloads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte(`MA-L,583653,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,58A15F,Texas Instruments,12500 TI Blvd Dallas TX US 75243`))
	}))
	defer server.Close()

	// Define a buffer to hold the downloaded database
	var buf bytes.Buffer

	// Download the database
	if err := oui.DownloadDatabase(&buf, server.URL); err != nil {
		t.Errorf("error returned from DownloadDatabase(): %v", err)
	}

	// Load the database
	db, err := oui.LoadDatabase(&buf)
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
}

// TestDownloadDatabaseServerError tests the case where the HTTP server returns
// an error. DownloadDatabase() should return the error below in this case.
// Error: "failed to download database file: 404 Not Found"
func TestDownloadDatabaseServerError(t *testing.T) {
	// Create a mock HTTP server for failed downloads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	// Define a buffer to hold the downloaded database
	var buf bytes.Buffer

	// Download the database
	if err := oui.DownloadDatabase(&buf, server.URL); err == nil {
		t.Errorf("expected error from DownloadDatabase(), got nil")
	}
}
