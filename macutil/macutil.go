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
package macutil

import (
	"fmt"
	"regexp"
)

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

// FindAllMacAddresses returns a list of MAC addresses found in the input string.
func FindAllMacAddresses(s string) ([]string, error) {
	// List of MAC addresses found in the input string
	var addresses []string

	// Define the MAC address systems to search for
	macSystem := []struct {
		groupSize int
		numGroups int
	}{
		{20, 2}, // IPoIB  : 00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01
		{8, 2},  // EUI-64 : 02:00:5e:10:00:00:00:01
		{6, 2},  // EUI-48 : 00:00:5e:00:53:01
		{10, 4}, // IPoIB  : 0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001
		{4, 4},  // EUI-64 : 0200.5e10.0000.0001
		{3, 4},  // EUI-48 : 0000.5e00.5301
	}

	// Extract MAC addresses from the input string
	input := s

	// Loop through the MAC address systems in order of most specific to least
	// specific. This is done to avoid false positives.
	for _, ms := range macSystem {
		// Extract MAC addresses from the output string. The output string is
		// updated with each iteration to remove the MAC addresses that were
		// found in the previous iteration.
		results, processedOutput, err := extractMacAddresses(input, ms.groupSize, ms.numGroups)
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
