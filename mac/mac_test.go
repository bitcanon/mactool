package mac

import (
	"testing"
)

// TestExtractMacAddresses tests the ExtractMacAddresses function.
func TestFindAllMacAddresses(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			input:    "00:00:5e:00:53:01",
			expected: []string{"00:00:5e:00:53:01"},
		},
		{
			input:    "00:00:5E:00:53:01",
			expected: []string{"00:00:5E:00:53:01"},
		},
		{
			input:    "02:00:5e:10:00:00:00:01",
			expected: []string{"02:00:5e:10:00:00:00:01"},
		},
		{
			input:    "00-00-5e-00-53-01",
			expected: []string{"00-00-5e-00-53-01"},
		},
		{
			input:    "02-00-5e-10-00-00-00-01",
			expected: []string{"02-00-5e-10-00-00-00-01"},
		},
		{
			input:    "0000.5e00.5301",
			expected: []string{"0000.5e00.5301"},
		},
		{
			input:    "0200.5e10.0000.0001",
			expected: []string{"0200.5e10.0000.0001"},
		},
		{
			input:    "0000-5e00-5301",
			expected: []string{"0000-5e00-5301"},
		},
		{
			input:    "0200-5e10-0000-0001",
			expected: []string{"0200-5e10-0000-0001"},
		},
		{
			input:    "MAC 1: 00:00:5E:00:53:01 and MAC 2: 0000.5E00.5301, done.",
			expected: []string{"00:00:5E:00:53:01", "0000.5E00.5301"},
		},
		{
			input:    "And a string without any addresses.",
			expected: nil,
		},
		{
			input:    "00005E-005301",
			expected: []string{"00005E-005301"},
		},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		// Find all MAC addresses in the input string
		macs, err := FindAllMacAddresses(tc.input)
		if err != nil {
			t.Errorf("error returned from FindAllMacAddresses(%q): %v", tc.input, err)
		}

		// Check the number of MAC addresses found
		if len(macs) != len(tc.expected) {
			t.Errorf("expected %d MAC addresses, got %d", len(tc.expected), len(macs))
		}

		// Compare the results to the expected values
		for i, mac := range macs {
			// Check the value of the MAC address
			if mac != tc.expected[i] {
				t.Errorf("expected %q, got %q", tc.expected[i], mac)
			}
		}
	}
}

// TestExtractOuiFromMac tests the ExtractOuiFromMac function.
func TestExtractOuiFromMac(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		input    string
		expected string
		expError error
	}{
		{"00:00:5e:00:53:01", "00005E", nil},
		{"00:00:5E:00:53:01", "00005E", nil},
		{"02:00:5e:10:00:00:00:01", "02005E", nil},
		{"00-00-5e-00-53-01", "00005E", nil},
		{"00005e-005301", "00005E", nil},
		{"02-00-5e-10-00-00-00-01", "02005E", nil},
		{"0000.5e00.5301", "00005E", nil},
		{"0200.5e10.0000.0001", "02005E", nil},
		{"0000-5e00-5301", "00005E", nil},
		{"0200-5e10-0000-0001", "02005E", nil},
		{"", "", ErrInvalidMacAddress},
		{"AB:CD", "", ErrInvalidMacAddress},
		{"NO:TA:MA:CA:DD:RE", "", ErrInvalidMacAddress},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		// Extract the OUI from the MAC address
		oui, err := ExtractOuiFromMac(tc.input)

		// Check for an expected error
		if err != tc.expError {
			t.Errorf("expected %v, got %v", tc.expError, err)
		}

		// Compare the results to the expected values
		if oui != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, oui)
		}
	}
}

// TestFindMacDelimiter tests the findMacDelimiter function.
func TestFindMacDelimiter(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name          string
		macAddress    string
		expectedDelim string
	}{
		{"WithColon", "00:11:22:33:44:55", ":"},
		{"WithHyphen", "00-11-22-33-44-55", "-"},
		{"WithDot", "0011.2233.4455", "."},
		{"WithInvalid", "00_11_22_33_44_55", ""},
		{"NoDelimiter", "001122334455", ""},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Find the delimiter in the MAC address
			actualDelim := findMacDelimiter(tc.macAddress)

			// Compare the results to the expected values
			if actualDelim != tc.expectedDelim {
				t.Errorf("expected '%s', but got '%s'", tc.expectedDelim, actualDelim)
			}
		})
	}
}

// TestCleanMacAddress tests the cleanMacAddress function.
func TestCleanMacAddress(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name       string
		macAddress string
		expected   string
	}{
		{"WithColon", "00:11:22:33:AA:55", "00112233AA55"},
		{"WithHyphen", "00-11-22-33-Aa-55", "00112233Aa55"},
		{"WithDot", "0011.2233.FF55", "00112233FF55"},
		{"WithDelimiters", "00:11-22.ABcd-efGH", "001122ABcdef"},
		{"AlphanumericOnly", "001122ABcdefGH", "001122ABcdef"},
		{"WithSpecialChars", "00#$%11@22", "001122"},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean the MAC address
			cleaned := cleanMacAddress(tc.macAddress)
			if cleaned != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, cleaned)
			}
		})
	}
}

// TestFormatWithDelimiters tests the formatWithDelimiters function.
func TestFormatWithDelimiters(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name          string
		macAddress    string
		delimiter     string
		groupSize     int
		expectedValue string
		expectedError error
	}{
		{"ColonDelimiter", "001122334455", ":", 2, "00:11:22:33:44:55", nil},
		{"HyphenDelimiter", "001122334455", "-", 2, "00-11-22-33-44-55", nil},
		{"DotDelimiter", "001122334455", ".", 2, "00.11.22.33.44.55", nil},
		{"NoDelimiter", "001122334455", "", 2, "001122334455", nil},
		{"GroupSizeIs2", "001122334455", ":", 2, "00:11:22:33:44:55", nil},
		{"GroupSizeIs3", "001122334455", ":", 3, "", ErrInvalidGroupSize},
		{"GroupSizeIs4", "001122334455", ".", 4, "0011.2233.4455", nil},
		{"InvalidMacLength", "00112233445", ".", 4, "", ErrInvalidMacAddressLength},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Format the MAC address
			formattedMac, err := formatWithDelimiters(tc.macAddress, tc.delimiter, tc.groupSize)

			// Compare the results to the expected values
			if err != tc.expectedError {
				t.Errorf("expected %v, but got %v", tc.expectedError, err)
			}

			if formattedMac != tc.expectedValue {
				t.Errorf("expected %s, but got %s", tc.expectedValue, formattedMac)
			}
		})
	}
}

// TestFormatMacAddress tests the FormatMacAddress function.
func TestFormatMacAddress(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		macAddress  string
		newFormat   MacFormat
		expectedMac string
		expectedErr error
	}{
		{"UpperColon", "001122AaBbCc", MacFormat{Upper, Colon, 2}, "00:11:22:AA:BB:CC", nil},
		{"LowerHyphen", "001122AaBbCc", MacFormat{Lower, Hyphen, 2}, "00-11-22-aa-bb-cc", nil},
		{"OriginalCaseDot", "001122AaBbCc", MacFormat{OriginalCase, Dot, 2}, "00.11.22.Aa.Bb.Cc", nil},
		{"UpperGroupSize2", "001122AaBbCc", MacFormat{Upper, Colon, 2}, "00:11:22:AA:BB:CC", nil},
		{"LowerGroupSize4", "001122AaBbCc", MacFormat{Lower, Dot, 4}, "0011.22aa.bbcc", nil},
		{"InvalidCaseOption", "001122AaBbCc", MacFormat{5, Hyphen, 2}, "", ErrInvalidCaseOption},
		{"OriginalCaseNoDelim", "001122AaBbCc", MacFormat{OriginalCase, None, 2}, "001122AaBbCc", nil},
		{"OriginalCaseInvalidDelim", "001122AaBbCc", MacFormat{OriginalCase, 5, 2}, "", ErrInvalidDelimiterOption},
		{"OriginalCaseOriginalDelim", "00:11:22:Aa:Bb:Cc", MacFormat{OriginalCase, OriginalDelim, 2}, "00:11:22:Aa:Bb:Cc", nil},
		{"OriginalCaseHyphenDelim", "00:11:22:Aa:Bb:Cc", MacFormat{OriginalCase, Hyphen, 2}, "00-11-22-Aa-Bb-Cc", nil},
		{"UpperHyphenDelim", "00:11:22:Aa:Bb:Cc", MacFormat{Upper, Hyphen, 2}, "00-11-22-AA-BB-CC", nil},
		{"UpperInvalidDelim", "001122AaBbCc", MacFormat{Upper, 11, 2}, "", ErrInvalidDelimiterOption},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Format the MAC address using the new format
			formattedMac, err := FormatMacAddress(tc.macAddress, tc.newFormat)

			// Compare the results to the expected values
			if tc.expectedErr == nil {
				// If no error is expected
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if formattedMac != tc.expectedMac {
					t.Errorf("expected %s, but got %s", tc.expectedMac, formattedMac)
				}
			} else {
				// If an error is expected
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
				}
			}
		})
	}
}

// TestGetGroupSize tests the getGroupSize function.
func TestGetGroupSize(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name        string
		macAddress  string
		expected    int
		expectedErr error
	}{
		{"WithColonGroupSizeIs2", "00:11:22:33:44:55", 2, nil},
		{"WithHyphenGroupSizeIs2", "00-11-22-33-44-55", 2, nil},
		{"WithPeriodGroupSizeIs2", "00.11.22.33.44.55", 2, nil},
		{"WithColonGroupSizeIs4", "0011:2233:4455", 4, nil},
		{"WithHyphenGroupSizeIs4", "0011-2233-4455", 4, nil},
		{"WithPeriodGroupSizeIs4", "0011.2233.4455", 4, nil},
		{"WithColonGroupSizeIs6", "001122:334455", 6, nil},
		{"WithHyphenGroupSizeIs6", "001122-334455", 6, nil},
		{"WithPeriodGroupSizeIs6", "001122.334455", 6, nil},
		{"WithPeriodGroupSizeIs3", "001.122:334-455", 3, nil},
		{"NoDelimiter", "001122334455", 0, ErrInvalidMacAddress},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the group size of the MAC address
			actual, err := getGroupSize(tc.macAddress)

			// Compare the results to the expected values
			if err != tc.expectedErr {
				t.Errorf("expected %v, but got %v", tc.expectedErr, err)
			}

			if actual != tc.expected {
				t.Errorf("expected %d, but got %d", tc.expected, actual)
			}
		})
	}
}
