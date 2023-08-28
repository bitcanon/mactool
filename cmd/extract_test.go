package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

// TestExtractAction tests the extractAction function
func TestExtractAction(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
		sortAsc  bool
		sortDesc bool
	}{
		{
			name:     "SingleLineInput",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			expected: "00:00:5e:00:53:01\n",
		},
		{
			name: "MultiLineInput",
			input: `First line of input with one MAC address 00:00:5e:00:53:01 in it.
			Second line of input with one MAC address 00-00-5E-00-53-02 in it.`,
			expected: "00:00:5e:00:53:01\n00-00-5E-00-53-02\n",
		},
		{
			name:     "SingleLineInputWithNoMacAddresses",
			input:    "First line of input with no MAC address in it.",
			expected: "",
		},
		{
			name: "MultiLineInputWithNoMacAddresses",
			input: `First line of input with no MAC address in it.
			Second line of input with no MAC address in it.`,
			expected: "",
		},
		{
			name:     "SingleLineInputWithMultipleMacAddresses",
			input:    "First MAC address 00:00:5e:00:53:01 and second MAC address 00-00-5E-00-53-02.",
			expected: "00:00:5e:00:53:01\n00-00-5E-00-53-02\n",
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: "",
		},
		{
			name:     "SingleLineInputWithSortAsc",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			sortAsc:  true,
			expected: "00:00:5e:00:53:01\n",
		},
		{
			name: "MultiLineInputWithSortAsc",
			input: `First line of input with one MAC address 22:22:5e:00:53:01 in it.
			Second line of input with one MAC address 11-11-5E-00-53-02 in it.`,
			sortAsc:  true,
			expected: "11-11-5E-00-53-02\n22:22:5e:00:53:01\n",
		},
		{
			name:     "SingleLineInputWithSortDesc",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			sortDesc: true,
			expected: "00:00:5e:00:53:01\n",
		},
		{
			name: "MultiLineInputWithSortDesc",
			input: `First line of input with one MAC address 11:11:5e:00:53:01 in it.
			Second line of input with one MAC address 22-22-5E-00-53-02 in it.`,
			sortDesc: true,
			expected: "22-22-5E-00-53-02\n11:11:5e:00:53:01\n",
		},
		{
			name:     "MultipleAddressesWithSortDesc",
			input:    `33:33:33:33:33:33 11:11:11:11:11:11 22:22:22:22:22:22 44:44:44:44:44:44 55:55:55:55:55:55`,
			sortDesc: true,
			expected: "55:55:55:55:55:55\n44:44:44:44:44:44\n33:33:33:33:33:33\n22:22:22:22:22:22\n11:11:11:11:11:11\n",
		},
		{
			name:     "MultipleAddressesWithSortAsc",
			input:    `33:33:33:33:33:33 11:11:11:11:11:11 22:22:22:22:22:22 44:44:44:44:44:44 55:55:55:55:55:55`,
			sortAsc:  true,
			expected: "11:11:11:11:11:11\n22:22:22:22:22:22\n33:33:33:33:33:33\n44:44:44:44:44:44\n55:55:55:55:55:55\n",
		},
	}

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Prepare a buffer to capture the output
			var output bytes.Buffer

			// Set the sort flags
			viper.Set("extract-sort-asc", test.sortAsc)
			viper.Set("extract-sort-desc", test.sortDesc)

			// Call the function to test
			err := extractAction(&output, test.input)

			// Check for errors
			if err != nil {
				t.Errorf("error returned from extractAction(): %v", err)
				return
			}

			// Compare the results to the expected values
			if output.String() != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, output.String())
			}
		})
	}
}
