/*
Copyright © 2023 Mikael Schultz <bitcanon@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package mac

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Errors returned by functions in this package
var ErrInvalidMacAddress = errors.New("invalid MAC address")
var ErrInvalidMacAddressLength = errors.New("invalid MAC address length; must be divisible by group size")
var ErrInvalidGroupSize = errors.New("invalid group size; must be 2, 4 or 6")
var ErrInvalidCaseOption = errors.New("invalid case option")
var ErrInvalidDelimiterOption = errors.New("invalid delimiter option")

// CaseOption is used to specify the case of the MAC address.
type CaseOption int

const (
	OriginalCase CaseOption = iota
	Upper
	Lower
)

// DelimiterOption is used to specify the delimiter used in the MAC address.
type DelimiterOption int

const (
	OriginalDelim DelimiterOption = iota
	Colon
	Hyphen
	Dot
	None
)

// GroupSizeOption is used to specify the number of characters in each group.
type GroupSizeOption int

const (
	OriginalGroupSize GroupSizeOption = iota
	GroupSizeTwo
	GroupSizeFour
	GroupSizeSix
)

// MacFormat is used to specify the format of the MAC address.
type MacFormat struct {
	Case      CaseOption
	Delimiter DelimiterOption
	GroupSize GroupSizeOption
}

// cleanMacAddress removes all non-alphanumeric characters from the MAC address.
func cleanMacAddress(macAddress string) string {
	// Use a regular expression to match non-alphanumeric characters
	reg := regexp.MustCompile("[^A-Fa-f0-9]")

	// Remove all non-alphanumeric characters from the MAC address
	return reg.ReplaceAllString(macAddress, "")
}

// extractMacAddresses extracts MAC addresses from the input string.
// The groupCount parameter is the number of groups of MAC
// address characters. The groupSize parameter is the number of
// characters in each group. Each group is separated by a colon,
// dash or period.
func extractMacAddresses(input string, groupCount int, groupSize int) ([]string, string, error) {
	// Regular expression pattern to match MAC addresses
	pattern := fmt.Sprintf(`((?:[\da-fA-F]{%d}[:\.-]){%d}[\da-fA-F]{%d})`,
		groupSize, groupCount-1, groupSize)

	// Compile the regular expression pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, input, err
	}

	// Save the MAC addresses found in the input string
	// and remove them from the input string
	addresses := re.FindAllString(input, -1)
	modifiedInput := re.ReplaceAllString(input, "")

	// Return the list of MAC addresses found in the input string
	// and the input string with the MAC addresses removed
	return addresses, modifiedInput, nil
}

// findMacDelimiter finds the delimiter used in the MAC address.
func findMacDelimiter(macAddress string) string {
	// Check for a colon, dash or period delimiter
	delimiters := []string{":", "-", "."}

	// Check for each delimiter
	for _, delimiter := range delimiters {
		// If the delimiter is found, return it
		if strings.Contains(macAddress, delimiter) {
			return delimiter
		}
	}

	// No delimiter was found
	return ""
}

// formatWithDelimiters formats the MAC address with the specified delimiter
// and group size. The group size is the number of characters in each group.
// The delimiter is the character used to separate each group.
func formatWithDelimiters(macAddress, delimiter string, groupSize int) (string, error) {
	// Validate groupSize
	if groupSize != 2 && groupSize != 4 && groupSize != 6 {
		return "", ErrInvalidGroupSize
	}

	// Ensure the MAC address contains only alphanumeric characters
	macAddress = cleanMacAddress(macAddress)

	// Validate length divisibility
	if len(macAddress)%groupSize != 0 {
		return "", ErrInvalidMacAddressLength
	}

	// Length of the MAC address
	macLength := len(macAddress)
	groupCount := macLength / groupSize

	// Split the MAC address into groups of characters
	// and join them with the delimiter
	groups := make([]string, groupCount)
	for i := 0; i < groupCount; i++ {
		groups[i] = macAddress[i*groupSize : i*groupSize+groupSize]
	}

	// Return the MAC address with the delimiter
	return strings.Join(groups, delimiter), nil
}

// FindAllMacAddresses returns a list of MAC addresses found in the input string.
func FindAllMacAddresses(s string) ([]string, error) {
	// List of MAC addresses found in the input string
	var addresses []string

	// Define the MAC address systems to search for
	macSystem := []struct {
		groupCount int
		groupSize  int
	}{
		{8, 2}, // EUI-64 : 02:00:5e:10:00:00:00:01
		{4, 4}, // EUI-64 : 0200.5e10.0000.0001
		{6, 2}, // EUI-48 : 00:00:5e:00:53:01
		{3, 4}, // EUI-48 : 0000.5e00.5301
		{2, 6}, // EUI-48 : 00005e-005301
	}

	// Extract MAC addresses from the input string
	input := s

	// Loop through the MAC address systems in order of most specific to least
	// specific. This is done to avoid false positives.
	for _, ms := range macSystem {
		// Extract MAC addresses from the output string. The output string is
		// updated with each iteration to remove the MAC addresses that were
		// found in the previous iteration.
		results, processedOutput, err := extractMacAddresses(input, ms.groupCount, ms.groupSize)
		if err != nil {
			return nil, err
		}

		// Append the results to the list of MAC addresses
		addresses = append(addresses, results...)

		// Update the output string to the processed output
		input = processedOutput
	}

	// Return the list of MAC addresses found in the input string
	return addresses, nil
}

// ExtractOuiFromMac extracts the OUI assignment from a MAC address.
func ExtractOuiFromMac(macAddress string) (string, error) {
	// Make sure the string is uppercase since the
	// assignment in the OUI database is uppercase
	macAddress = strings.ToUpper(macAddress)

	// Define a regular expression to remove all non-alphanumeric
	// characters from the MAC address
	reg := regexp.MustCompile("[^a-fA-F0-9]+")

	// Remove all non-alphanumeric characters from the MAC address
	macAddress = reg.ReplaceAllString(macAddress, "")

	// Make sure the MAC address is at least 12 characters long
	// and return the first 3 bytes (6 hexadecimal characters) as a string
	if len(macAddress) >= 12 {
		assignment := macAddress[0:6]
		return assignment, nil
	} else {
		return "", ErrInvalidMacAddress
	}
}

// FormatMacAddress formats the MAC address with the specified case, delimiter
// and group size. The group size is the number of characters in each group.
// The delimiter is the character used to separate each group.
// Example: 00:00:5E:00:53:01 (Case: Upper, Delimiter: Colon, GroupSize: 2)
func FormatMacAddress(macAddress string, newFormat MacFormat) (string, error) {
	// Validate the case option
	switch newFormat.Case {
	case Upper:
		macAddress = strings.ToUpper(macAddress)
	case Lower:
		macAddress = strings.ToLower(macAddress)
	case OriginalCase:
		// Keep the original case
	default:
		// Invalid case option specified
		return "", ErrInvalidCaseOption
	}

	// Validate the delimiter option
	delimiterStr := ""
	switch newFormat.Delimiter {
	case Colon:
		delimiterStr = ":"
	case Hyphen:
		delimiterStr = "-"
	case Dot:
		delimiterStr = "."
	case None:
		delimiterStr = ""
	case OriginalDelim:
		// Keep the original delimiter
		delimiterStr = findMacDelimiter(macAddress)
	default:
		// Invalid delimiter option specified
		return "", ErrInvalidDelimiterOption
	}

	// Validate the group size option
	groupSize := 0
	switch newFormat.GroupSize {
	case GroupSizeTwo:
		groupSize = 2
	case GroupSizeFour:
		groupSize = 4
	case GroupSizeSix:
		groupSize = 6
	case OriginalGroupSize:
		groupSize, _ = GetGroupSize(macAddress)
	}

	// Format the MAC address
	mac, err := formatWithDelimiters(macAddress, delimiterStr, groupSize)
	if err != nil {
		return "", err
	}

	// Return the formatted MAC address
	return mac, nil
}

// getGroupSize calculates the group size of the MAC address.
// The group size is the number of characters in each group.
// A delimiting character separates each group, and can be any
// character except alphanumeric characters.
func GetGroupSize(macAddress string) (int, error) {
	// Remove all non-alphanumeric characters from the MAC address
	strippedMAC := regexp.MustCompile(`[^0-9a-zA-Z]`).ReplaceAllString(macAddress, "")
	strippedLen := len(strippedMAC)

	// Calculate the number of delimiters removed from the MAC address
	difference := len(macAddress) - strippedLen

	// If the difference is 0, the MAC address is invalid
	if difference == 0 {
		return 0, ErrInvalidMacAddress
	}

	// Calculate the group size
	groupSize := strippedLen / (difference + 1)

	// If group size is not 2, 4 or 6, the MAC address is invalid
	if groupSize != 2 && groupSize != 4 && groupSize != 6 {
		return 0, ErrInvalidMacAddress
	}

	// Return the group size and no error
	return groupSize, nil
}
