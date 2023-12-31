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

// TestLen tests the Len function, which returns the number of entries in the database.
func TestLen(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,583653,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,58A15F,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Verify that the database was loaded correctly
	if db.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", db.Len())
	}
}

// TestLess tests the Less function use for sorting.
func TestLess(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,111111,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,222222,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Verify that the database was loaded correctly
	if db.Less(0, 1) != true {
		t.Errorf("expected true, got %v", db.Less(0, 1))
	}
}

// TestSwap tests the case where the database is loaded successfully.
// The entries at the specified indices should be swapped.
func TestSwap(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,111111,"Apple, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,222222,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Verify that the database was loaded correctly
	db.Swap(0, 1)

	// Verify that the database was loaded correctly
	if db.Entries[0].Assignment != "222222" {
		t.Errorf("expected 222222, got %s", db.Entries[0].Assignment)
	}
}

// TestFindAllVendors tests the FindAllVendors function with various inputs.
func TestFindAllVendors(t *testing.T) {
	// Create a test CSV database
	csvData := `MA-L,111111,"Banana, Inc.",1 Infinite Noob Cupertino CA US 22222
MA-L,111222,Texas Instruments,11111 TI Blvd Dallas TX US 75243
MA-L,222222,Texas Instruments,12500 TI Blvd Dallas TX US 75243`

	// Load the test CSV database
	db, err := oui.LoadDatabase(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("error returned from LoadDatabase(): %v", err)
	}

	// Setup test cases
	testCases := []struct {
		input         string
		filterOptions oui.FilterOptions
		expected      int
	}{
		{input: "111", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: false}, expected: 2},
		{input: "222", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: false}, expected: 3},
		{input: "333", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: false}, expected: 0},
		{input: "111", filterOptions: oui.FilterOptions{Assignment: true, Organization: false, Address: false}, expected: 2},
		{input: "1111", filterOptions: oui.FilterOptions{Assignment: true, Organization: false, Address: false}, expected: 1},
		{input: "333", filterOptions: oui.FilterOptions{Assignment: true, Organization: false, Address: false}, expected: 0},
		{input: "Banana", filterOptions: oui.FilterOptions{Assignment: false, Organization: true, Address: false}, expected: 1},
		{input: "banana", filterOptions: oui.FilterOptions{Assignment: false, Organization: true, Address: false}, expected: 1},
		{input: "in", filterOptions: oui.FilterOptions{Assignment: false, Organization: true, Address: false}, expected: 3},
		{input: "1111", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: true}, expected: 1},
		{input: "US", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: true}, expected: 3},
		{input: "222", filterOptions: oui.FilterOptions{Assignment: false, Organization: false, Address: true}, expected: 1},
	}

	// Loop through the test cases
	for _, testCase := range testCases {
		// Find all vendors matching the input string
		vendors, err := db.FindAllVendors(testCase.input, testCase.filterOptions)
		if err != nil {
			t.Errorf("error returned from FindAllVendors(): %v", err)
		}

		// Verify that the database was loaded correctly
		if len(vendors.Entries) != testCase.expected {
			t.Errorf("expected %d entries, got %d", testCase.expected, len(vendors.Entries))
		}
	}
}

// TestOuiContains tests the Contains function of the Oui type.
func TestOuiContains(t *testing.T) {
	// Create a test CSV database
	entry := oui.Oui{
		Assignment:   "1A2B3C",
		Organization: "Banana, Inc.",
		Address:      "1 Infinite Fruity Loop CA US 12014",
	}

	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:  "FullAssignment",
			input: "1A2B3C", expected: true,
		},
		{
			name:  "PartialAssignment",
			input: "1A2B3", expected: true,
		},
		{
			name:  "PartialAssignmentLowercase",
			input: "1a2b3", expected: true,
		},
		{
			name:  "FullOrganization",
			input: "Banana, Inc.", expected: true,
		},
		{
			name:  "PartialOrganizationLowercase",
			input: "banana", expected: true,
		},
		{
			name:  "FullAddress",
			input: "1 Infinite Fruity Loop CA US 12014", expected: true,
		},
		{
			name:  "PartialAddressUppercase",
			input: "LOOP", expected: true,
		},
		{
			name:  "PartialAddressLowercase",
			input: "fruit", expected: true,
		},
		{
			name:  "NotFound",
			input: "111222", expected: false},
	}

	// Loop through the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Verify that the database was loaded correctly
			if entry.Contains(testCase.input) != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, entry.Contains(testCase.input))
			}
		})
	}
}
