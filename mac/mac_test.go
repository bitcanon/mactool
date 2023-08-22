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
		{"", nil},
		{"00:00:5e:00:53:01", []string{"00:00:5e:00:53:01"}},
		{"00:00:5E:00:53:01", []string{"00:00:5E:00:53:01"}},
		{"02:00:5e:10:00:00:00:01", []string{"02:00:5e:10:00:00:00:01"}},
		{
			"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01",
			[]string{"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01"},
		},
		{"00-00-5e-00-53-01", []string{"00-00-5e-00-53-01"}},
		{"02-00-5e-10-00-00-00-01", []string{"02-00-5e-10-00-00-00-01"}},
		{
			"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01",
			[]string{"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01"},
		},
		{"0000.5e00.5301", []string{"0000.5e00.5301"}},
		{"0200.5e10.0000.0001", []string{"0200.5e10.0000.0001"}},
		{
			"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001",
			[]string{"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001"},
		},
		{"0000-5e00-5301", []string{"0000-5e00-5301"}},
		{"0200-5e10-0000-0001", []string{"0200-5e10-0000-0001"}},
		{
			"0000-0000-fe80-0000-0000-0000-0200-5e10-0000-0001",
			[]string{"0000-0000-fe80-0000-0000-0000-0200-5e10-0000-0001"},
		},
		{
			"MAC 1: 00:00:5E:00:53:01 and MAC 2: 0000.5E00.5301, done.",
			[]string{"00:00:5E:00:53:01", "0000.5E00.5301"},
		},
		{"And a string without any addresses.", nil},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		// Find all MAC addresses in the input string
		macs, err := FindAllMacAddresses(tc.input)
		if err != nil {
			t.Errorf("error returned from FindAllMacAddresses(%q): %v", tc.input, err)
		}

		// Compare the results to the expected values
		for i, mac := range macs {
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
