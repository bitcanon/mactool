package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/bitcanon/mactool/mac"
	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/viper"
)

// TestFormatAction tests the formatAction function.
func TestFormatAction(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name     string
		input    string
		expected string
		format   mac.MacFormat
	}{
		{
			name: "FormatUnmodified",
			input: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.OriginalDelim,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatCleanMacUnmodified",
			input: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.OriginalDelim,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatHyphenUpper",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.Hyphen,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatColonLower",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00:1a:2b:3c:4d:5e
				And another one: aa:bb:cc:dd:ee:ff
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.Lower,
				Delimiter: mac.Colon,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatDotOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-Dd-Ee-Ff
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00.1A.2B.3C.4D.5E
				And another one: AA.BB.CC.Dd.Ee.Ff
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Dot,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatNoDelimOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.None,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatColonOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA:BB:CC:DD:EE:FF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Colon,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatHyphenOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Hyphen,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatDotLower",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00.1a.2b.3c.4d.5e
				And another one: aa.bb.cc.dd.ee.ff
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.Lower,
				Delimiter: mac.Dot,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatNoDelimUpper",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.None,
				GroupSize: mac.OriginalGroupSize,
			},
		},
		{
			name: "FormatColonUpperGroupSize4",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: 00-1A-2B-3C-4D-5E
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A:2B3C:4D5E
				And another one: 001A:2B3C:4D5E
				No MAC here: 123456` + "\n",
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.Colon,
				GroupSize: mac.GroupSizeFour,
			},
		},
	}

	// Loop through the test cases and run each test
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Set up the input and output buffers
			var output strings.Builder

			// Call the formatAction function
			err := formatAction(&output, test.format, test.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Compare the output
			actualOutput := output.String()
			if actualOutput != test.expected {
				t.Errorf("expected '%s'\n\nbut got '%s'", test.expected, actualOutput)
			}
		})
	}
}

// TestFormatActionOutput tests the formatAction function
// with the output and append options set and unset.
func TestFormatActionOutput(t *testing.T) {
	// Get a tempfile for the output
	outputFile, err := os.CreateTemp("", "extract_output.txt")
	if err != nil {
		t.Errorf("error creating tempfile: %v", err)
		return
	}
	defer os.Remove(outputFile.Name())

	// Set up test cases
	testCases := []struct {
		name       string
		input      string
		expected   string
		format     mac.MacFormat
		outputFile string
		append     bool
	}{
		{
			name:       "EmptyInput",
			input:      "",
			expected:   "",
			format:     mac.MacFormat{Case: mac.Lower, Delimiter: mac.Colon, GroupSize: 2},
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name:       "SingleLineInput",
			input:      "Single line of input with one MAC address 00-00-5e-00-53-01 in it.",
			expected:   "Single line of input with one MAC address 00:00:5e:00:53:01 in it." + "\n",
			format:     mac.MacFormat{Case: mac.Lower, Delimiter: mac.Colon, GroupSize: mac.OriginalGroupSize},
			outputFile: outputFile.Name(),
			append:     false,
		},
		{
			name: "MultiLineInput",
			input: `First line of input with one MAC address 00-00-5e-00-53-01 in it.
Second line of input with one MAC address 00-00-5E-00-53-02 in it.`,
			expected: `Single line of input with one MAC address 00:00:5e:00:53:01 in it.
First line of input with one MAC address 00:00:5e:00:53:01 in it.
Second line of input with one MAC address 00:00:5e:00:53:02 in it.` + "\n",
			format:     mac.MacFormat{Case: mac.Lower, Delimiter: mac.Colon, GroupSize: mac.OriginalGroupSize},
			outputFile: outputFile.Name(),
			append:     true,
		},
	}

	// Loop through the test cases and run each test
	for _, test := range testCases {
		// Set the output file
		viper.Set("format.output", test.outputFile)
		viper.Set("format.append", test.append)

		// Get the output stream
		outStream, err := utils.GetOutputStream(test.outputFile, test.append)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		defer outStream.Close()

		// Call the function to test
		err = formatAction(outStream, test.format, test.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Read the output file using os.ReadFile()
		actual, err := os.ReadFile(test.outputFile)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Compare the output
		if string(actual) != test.expected {
			t.Errorf("expected '%s'\n\nbut got '%s'", test.expected, string(actual))
		}
	}
}
