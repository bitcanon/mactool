package mac_test

import (
	"testing"

	"github.com/bitcanon/mactool/mac"
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
		macs, err := mac.FindAllMacAddresses(tc.input)
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
		{"", "", mac.ErrInvalidMacAddress},
		{"AB:CD", "", mac.ErrInvalidMacAddress},
		{"NO:TA:MA:CA:DD:RE", "", mac.ErrInvalidMacAddress},
	}

	// Loop through the test cases
	for _, tc := range testCases {
		// Extract the OUI from the MAC address
		oui, err := mac.ExtractOuiFromMac(tc.input)

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
