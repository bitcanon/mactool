package cmd

import (
	"strings"
	"testing"

	"github.com/bitcanon/mactool/mac"
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
			name: "FormatHyphenUpper",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.Hyphen,
				GroupSize: 2,
			},
		},
		{
			name: "FormatColonLower",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00:1a:2b:3c:4d:5e
				And another one: aa:bb:cc:dd:ee:ff
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.Lower,
				Delimiter: mac.Colon,
				GroupSize: 2,
			},
		},
		{
			name: "FormatDotOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-Dd-Ee-Ff
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00.1A.2B.3C.4D.5E
				And another one: AA.BB.CC.Dd.Ee.Ff
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Dot,
				GroupSize: 2,
			},
		},
		{
			name: "FormatNoDelimOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.None,
				GroupSize: 2,
			},
		},
		{
			name: "FormatColonOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA:BB:CC:DD:EE:FF
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Colon,
				GroupSize: 2,
			},
		},
		{
			name: "FormatHyphenOriginalCase",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00-1A-2B-3C-4D-5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.OriginalCase,
				Delimiter: mac.Hyphen,
				GroupSize: 2,
			},
		},
		{
			name: "FormatDotLower",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 00.1a.2b.3c.4d.5e
				And another one: aa.bb.cc.dd.ee.ff
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.Lower,
				Delimiter: mac.Dot,
				GroupSize: 2,
			},
		},
		{
			name: "FormatNoDelimUpper",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: AA-BB-CC-DD-EE-FF
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A2B3C4D5E
				And another one: AABBCCDDEEFF
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.None,
				GroupSize: 2,
			},
		},
		{
			name: "FormatColonUpperGroupSize4",
			input: `Here is a MAC address: 00:1A:2B:3C:4D:5E
				And another one: 00-1A-2B-3C-4D-5E
				No MAC here: 123456`,
			expected: `Here is a MAC address: 001A:2B3C:4D5E
				And another one: 001A:2B3C:4D5E
				No MAC here: 123456`,
			format: mac.MacFormat{
				Case:      mac.Upper,
				Delimiter: mac.Colon,
				GroupSize: 4,
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
