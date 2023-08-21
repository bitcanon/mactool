package cmd

import (
	"bytes"
	"testing"
)

// TestExtractAction tests the extractAction function
func TestExtractAction(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    string
		expected string
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
	}

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Prepare a buffer to capture the output
			var output bytes.Buffer

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
