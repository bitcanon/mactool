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
package utils

import "os"

// GetOutputStream returns an output stream for the specified filename
// If the filename is empty, the standard output stream is returned
// If the append flag is true, the file is opened for appending
// If the append flag is false, the file is opened for writing
func GetOutputStream(filename string, append bool) (*os.File, error) {
	// fileMode controls how the file is opened
	var fileMode int

	// If no filename is specified, use standard output
	if filename == "" {
		return os.Stdout, nil
	}

	if append {
		// If append is true, append to the file
		fileMode = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else {
		// If append is false, overwrite the file
		fileMode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	// Open the file for writing using the specified file mode
	outStream, err := os.OpenFile(filename, fileMode, 0644)
	if err != nil {
		return nil, err
	}

	// Return the output stream
	return outStream, nil
}
