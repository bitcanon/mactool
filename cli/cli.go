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
package cli

import (
	"bufio"
	"io"
	"os"
)

// processStdin reads all data from standard input
// and returns the input as a string
func ProcessStdin() (string, error) {
	// Read all data from standard input
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	// Return the input string on success
	return string(input), nil
}

// processInteractiveInput processes the interactive input
// and extracts MAC addresses from the input string
func ProcessInteractiveInput() (string, error) {
	// A string for the user input
	var input string

	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Read each line from standard input as the user types.
	// The loop will exit when the user presses Ctrl+D (Unix)
	// or Ctrl+Z (Windows).
	for scanner.Scan() {
		input = scanner.Text()
	}

	// Check for errors that may have occurred while reading
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Return the input string on success
	return input, nil
}
