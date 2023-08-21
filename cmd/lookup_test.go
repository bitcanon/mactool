package cmd

import (
	"strings"
	"testing"
)

func TestLookupAction(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SingleLineInput",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n",
		},
		{
			name: "MultiLineInput",
			input: `First line of input with one MAC address 00:00:5e:00:53:01 in it.
			Second line of input with one MAC address 00-00-5E-00-53-02 in it.`,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n00-00-5E-00-53-02 (Banana, Inc.)\n",
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
			input:    "First MAC address 00:00:5e:00:53:01 and second MAC address 12-3A-BC-00-53-02.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12-3A-BC-00-53-02 (Swede Instruments)\n",
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: "",
		},
	}

	// Create a test CSV database
	csvData := `Registry,Assignment,Organization Name,Organization Address
MA-L,00005E,"Banana, Inc.",1 Infinite Loop Cupocoffee CA US 12345
MA-L,123ABC,Swede Instruments,Storgatan 1 Stockholm SE 12345`

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a reader for the test CSV database
			reader := strings.NewReader(csvData)

			// Prepare a buffer to capture the output
			var output strings.Builder

			// Call the function to test
			err := lookupAction(&output, reader, test.input)
			if err != nil {
				t.Errorf("error returned from lookupAction(): %v", err)
				return
			}

			// Check the output
			if output.String() != test.expected {
				t.Errorf("lookupAction() output = %q, want %q", output.String(), test.expected)
			}
		})
	}

}
