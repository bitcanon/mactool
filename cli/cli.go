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
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// ProcessStdin reads all data from standard input
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

// ProcessInteractiveInput processes the interactive input
// and extracts MAC addresses from the input string
func ProcessInteractiveInput() (string, error) {
	// A string for the user input
	var input string

	// Create a scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Tell the user how to finish the input
	// based on the operating system
	eofKeys := "CTRL+D"
	if runtime.GOOS == "windows" {
		eofKeys = "CTRL+Z"
	}
	fmt.Fprintf(os.Stderr, "Please enter the input text. Press %s to finish.\n", eofKeys)

	// Read each line from standard input as the user types
	for scanner.Scan() {
		input += fmt.Sprintf("%s\n", scanner.Text())
	}

	// Remove the trailing newline character
	input = strings.TrimRight(input, "\n")

	// Check for errors that may have occurred while reading
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Return the input string on success
	return input, nil
}

// ProcessFile reads all data from the specified file
// and returns the input as a string
func ProcessFile(filename string) (string, error) {
	// Open the input file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a scanner to read from the file
	scanner := bufio.NewScanner(file)

	// A string for the file contents
	var input string

	// Read each line from the file
	for scanner.Scan() {
		input += fmt.Sprintf("%s\n", scanner.Text())
	}

	// Remove the trailing newline character
	input = strings.TrimRight(input, "\n")

	// Check for errors that may have occurred while reading
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Return the input string on success
	return input, nil
}
