package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/viper"
)

// TestLookupAction tests the lookupAction function
func TestLookupAction(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name     string
		input    string
		expected string
		sortAsc  bool
		sortDesc bool
		suppress bool
	}{
		{
			name:     "SingleLineInput",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n",
			suppress: false,
		},
		{
			name: "MultiLineInput",
			input: `First line of input with one MAC address 00:00:5e:00:53:01 in it.
			Second line of input with one MAC address 00-00-5E-00-53-02 in it.`,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n00-00-5E-00-53-02 (Banana, Inc.)\n",
			suppress: false,
		},
		{
			name:     "SingleLineInputWithNoMacAddresses",
			input:    "First line of input with no MAC address in it.",
			expected: "",
			suppress: false,
		},
		{
			name: "MultiLineInputWithNoMacAddresses",
			input: `First line of input with no MAC address in it.
			Second line of input with no MAC address in it.`,
			expected: "",
			suppress: false,
		},
		{
			name:     "SingleLineInputWithMultipleMacAddresses",
			input:    "First MAC address 00:00:5e:00:53:01 and second MAC address 12-3A-BC-00-53-02.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12-3A-BC-00-53-02 (Swede Instruments)\n",
			suppress: false,
		},
		{
			name:     "SingleLineInputWithNoVendorFound",
			input:    "First line of input with one MAC address 00:11:22:00:53:01 in it.",
			expected: "00:11:22:00:53:01\n",
			suppress: false,
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: "",
			suppress: false,
		},
		{
			name:     "InputWithAllMacAddressesSuppressed",
			input:    "First MAC address 99:99:99:00:53:01 and second MAC address 99-99-99-00-53-02.",
			expected: "",
			suppress: true,
		},
		{
			name:     "InputWithSomeMacAddressesSuppressed",
			input:    "First MAC address 00:00:5e:00:53:01 and second MAC address 99-99-99-00-53-02.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n",
			suppress: true,
		},
		{
			name:     "InputWithNoMacAddressesSuppressed",
			input:    "First MAC address 00:00:5e:00:53:01 and second MAC address 12-3A-BC-00-53-02.",
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12-3A-BC-00-53-02 (Swede Instruments)\n",
			suppress: true,
		},
		{
			name:     "SingleLineInputWithSortAsc",
			input:    "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			sortAsc:  true,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n",
		},
		{
			name: "MultiLineInputWithSortAsc",
			input: `First line of input with one MAC address 12:3A:BC:00:53:01 in it.
			Second line of input with one MAC address 00:00:5e:00:53:01 in it.`,
			sortAsc:  true,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12:3A:BC:00:53:01 (Swede Instruments)\n",
		},
		{
			name: "MultiLineInputWithSortAscAndDesc",
			input: `First line of input with one MAC address 12:3A:BC:00:53:01 in it.
			Second line of input with one MAC address 00:00:5e:00:53:01 in it.`,
			sortAsc:  true,
			sortDesc: true,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12:3A:BC:00:53:01 (Swede Instruments)\n",
		},
		{
			name: "MultiLineInputWithSortDesc",
			input: `First line of input with one MAC address 00:00:5e:00:53:01 in it.
			Second line of input with one MAC address 12-3A-BC-00-53-02 in it.`,
			sortDesc: true,
			expected: "12-3A-BC-00-53-02 (Swede Instruments)\n00:00:5e:00:53:01 (Banana, Inc.)\n",
		},
		{
			name:     "SingleLineWithSortAscAndSuppress",
			input:    "MAC1: 12:3A:BC:00:53:01, MAC2: 00:00:5e:00:53:01 and MAC3: 99-99-99-00-53-02.",
			sortAsc:  true,
			suppress: true,
			expected: "00:00:5e:00:53:01 (Banana, Inc.)\n12:3A:BC:00:53:01 (Swede Instruments)\n",
		},
		{
			name:     "SingleLineWithSortDescAndSuppress",
			input:    "MAC1: 00:00:5e:00:53:01, MAC2: 99-99-99-00-53-02 and MAC3: 12:3A:BC:00:53:01.",
			sortDesc: true,
			suppress: true,
			expected: "12:3A:BC:00:53:01 (Swede Instruments)\n00:00:5e:00:53:01 (Banana, Inc.)\n",
		},
	}

	// Create a test CSV database, in memory, to be used by the test cases
	csvData := `Registry,Assignment,Organization Name,Organization Address
MA-L,00005E,"Banana, Inc.",1 Infinite Loop Cupocoffee CA US 12345
MA-L,123ABC,Swede Instruments,Storgatan 1 Stockholm SE 12345`

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a reader for the test CSV database
			reader := strings.NewReader(csvData)

			// Set up viper with the suppress-unmatched flag
			viper.Set("lookup.suppress-unmatched", test.suppress)

			// Set the sort flags
			viper.Set("lookup.sort-asc", test.sortAsc)
			viper.Set("lookup.sort-desc", test.sortDesc)

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

func TestLookupActionOutput(t *testing.T) {
	outputFile, err := os.CreateTemp("", "lookup_output.txt")
	if err != nil {
		t.Errorf("error creating tempfile: %v", err)
		return
	}
	defer os.Remove(outputFile.Name())

	// Setup test cases
	testCases := []struct {
		name       string
		input      string
		expected   string
		outputFile string
		append     bool
	}{
		{
			name:       "SingleLineInput",
			input:      "First line of input with one MAC address 00:00:5e:00:53:01 in it.",
			expected:   "00:00:5e:00:53:01 (Banana, Inc.)\n",
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name: "MultiLineInput",
			input: `First line of input with one MAC address 00:00:5e:00:53:01 in it.
			Second line of input with one MAC address 00-00-5E-00-53-02 in it.`,
			expected:   "00:00:5e:00:53:01 (Banana, Inc.)\n00-00-5E-00-53-02 (Banana, Inc.)\n",
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name:       "SingleLineInputWithNoMacAddresses",
			input:      "First line of input with no MAC address in it.",
			expected:   "",
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name:       "SingleLineInputWithMultipleMacAddresses",
			input:      `First MAC address 00:00:5e:11:11:11 and second MAC address 12-3A-BC-00-53-02 in it.`,
			expected:   "00:00:5e:11:11:11 (Banana, Inc.)\n12-3A-BC-00-53-02 (Swede Instruments)\n",
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name:       "SingleLineInputWithAppend",
			input:      "One MAC address 00:00:5e:22:22:22 in this line.",
			expected:   "00:00:5e:11:11:11 (Banana, Inc.)\n12-3A-BC-00-53-02 (Swede Instruments)\n00:00:5e:22:22:22 (Banana, Inc.)\n",
			outputFile: outputFile.Name(),
			append:     true,
		},
	}

	// Create a test CSV database, in memory, to be used by the test cases
	csvData := `Registry,Assignment,Organization Name,Organization Address
MA-L,00005E,"Banana, Inc.",1 Infinite Loop Cupocoffee CA US 12345
MA-L,123ABC,Swede Instruments,Storgatan 1 Stockholm SE 12345`

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a reader for the test CSV database
			reader := strings.NewReader(csvData)

			// Disable sorting for this test
			viper.Set("lookup.sort-asc", false)
			viper.Set("lookup.sort-desc", false)

			// Set the output file
			viper.Set("lookup.output-file", test.outputFile)
			viper.Set("lookup.append", test.append)

			// Get the output stream
			outStream, err := utils.GetOutputStream(test.outputFile, test.append)
			if err != nil {
				t.Errorf("error getting output stream: %v", err)
				return
			}
			defer outStream.Close()

			// Call the function to test
			err = lookupAction(outStream, reader, test.input)
			if err != nil {
				t.Errorf("error returned from lookupAction(): %v", err)
				return
			}

			// Read the output file using os.ReadFile()
			actual, err := os.ReadFile(test.outputFile)
			if err != nil {
				t.Errorf("error reading output file: %v", err)
				return
			}

			// Compare the results to the expected values
			if string(actual) != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, string(actual))
			}
		})
	}
}
