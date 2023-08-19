package cli_test

import (
	"os"
	"testing"

	"github.com/bitcanon/mactool/cli"
)

// TestProcessStdin tests the ProcessStdin function
// by redirecting stdin to a pipe and writing test
// data to the pipe to simulate user input
func TestProcessStdin(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:     "ValidInput",
			input:    "Example input from stdin",
			expected: "Example input from stdin",
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Redirect stdin for the test
			originalStdin := os.Stdin

			// Restore stdin when the test is done
			defer func() { os.Stdin = originalStdin }()

			// Create a pipe for stdin and write test data to it
			r, w, _ := os.Pipe()
			os.Stdin = r
			w.WriteString(test.input)
			w.Close()
			defer r.Close()

			// Process stdin using the cli package
			output, err := cli.ProcessStdin()
			if err != nil {
				t.Errorf("error returned from ProcessStdin(): %v", err)
				return
			}

			// Compare the results to the expected values
			if output != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, output)
			}
		})
	}
}

func TestProcessInteractiveInput(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:     "ValidSingleLineInput",
			input:    `First line of input from interactive user input`,
			expected: `First line of input from interactive user input`,
		},
		{
			name: "ValidMultiLineInput",
			input: `First line of input from interactive user input
Second line of input from interactive user input`,
			expected: `First line of input from interactive user input
Second line of input from interactive user input`,
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Redirect stdin for the test
			originalStdin := os.Stdin

			// Restore stdin when the test is done
			defer func() { os.Stdin = originalStdin }()

			// Create a pipe to simulate user input
			r, w, err := os.Pipe()
			if err != nil {
				t.Errorf("error returned from os.Pipe(): %v", err)
				return
			}

			// Redirect stdin to the pipe and close it when done
			os.Stdin = r
			defer r.Close()

			// Write test data to the pipe in a separate goroutine
			// to simulate user input in interactive mode (asynchronously)
			go func() {
				defer w.Close()
				w.Write([]byte(test.input))
			}()

			// Process stdin using the cli package
			output, err := cli.ProcessInteractiveInput()
			if err != nil {
				t.Errorf("error returned from ProcessInteractiveInput(): %v", err)
				return
			}

			// Compare the results to the expected values
			if output != test.expected {
				t.Errorf("expected %q, but got %q", test.expected, output)
			}
		})
	}
}
