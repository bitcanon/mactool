package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

// TestLookupVendorAction tests the lookupVendorAction function
func TestLookupVendorAction(t *testing.T) {
	// Create the test CSV database
	csvData := `MA-L,111111,"Banana, Inc.",1 Infinite Loop Cupertino CA US 12514
MA-L,222222,"Banana, Inc.",1 Infinite Loop Cupertino CA US 95014
MA-L,ABCDEF,Swede Instruments CA,12300 TI Blvd Dallas TX US 75243
MA-L,ABCABC,Sweet Instruments,12500 TI Blvd Dallas TX US 75243`

	// Setup test cases
	testCases := []struct {
		name            string
		input           string
		expected        string
		setAssignment   bool
		setOrganization bool
		setAddress      bool
	}{
		{
			name:     "FindInAllColumns",
			input:    "ca",
			expected: "111111 Banana, Inc.\n222222 Banana, Inc.\nABCDEF Swede Instruments CA\nABCABC Sweet Instruments\n",
		},
		{
			name:     "FindInAllColumnsNoMatch",
			input:    "nothere",
			expected: "",
		},
		{
			name:          "FindInAssignmentColumn",
			input:         "ABC",
			expected:      "ABCDEF Swede Instruments CA\nABCABC Sweet Instruments\n",
			setAssignment: true,
		},
		{
			name:          "FindInAssignmentColumnLowercase",
			input:         "abc",
			expected:      "ABCDEF Swede Instruments CA\nABCABC Sweet Instruments\n",
			setAssignment: true,
		},
		{
			name:            "FindInOrganizationColumn",
			input:           "in",
			expected:        "111111 Banana, Inc.\n222222 Banana, Inc.\nABCDEF Swede Instruments CA\nABCABC Sweet Instruments\n",
			setOrganization: true,
		},
		{
			name:            "FindInOrganizationColumn2",
			input:           "CA",
			expected:        "ABCDEF Swede Instruments CA\n",
			setOrganization: true,
		},
		{
			name:       "FindInAddressColumn",
			input:      "CA",
			expected:   "111111 Banana, Inc.\n222222 Banana, Inc.\n",
			setAddress: true,
		},
		{
			name:       "FindInAddressColumn2",
			input:      "us",
			expected:   "111111 Banana, Inc.\n222222 Banana, Inc.\nABCDEF Swede Instruments CA\nABCABC Sweet Instruments\n",
			setAddress: true,
		},
		{
			name:       "FindInAddressColumn3",
			input:      "111111",
			expected:   "",
			setAddress: true,
		},
	}

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a buffer to hold the output
			var buf bytes.Buffer

			// Create the CSV io reader
			csvReader := bytes.NewReader([]byte(csvData))

			// Set the filter options
			viper.Set("lookup-vendor.assignment", test.setAssignment)
			viper.Set("lookup-vendor.organization", test.setOrganization)
			viper.Set("lookup-vendor.address", test.setAddress)

			// Run the test
			err := lookupVendorAction(&buf, csvReader, test.input)
			if err != nil {
				t.Errorf("error returned from lookupVendorAction(): %v", err)
			}

			// Verify that the output matches the expected output
			if buf.String() != test.expected {
				t.Errorf("expected '%s', got '%s'", test.expected, buf.String())
			}
		})
	}

}
