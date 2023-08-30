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

// TestProcessInteractiveInput tests the ProcessInteractiveInput function
// by redirecting stdin to a pipe and writing test data to the pipe to
// simulate user input
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

// TestProcessFile tests the ProcessFile function
// by creating a temporary file and writing test
// data to it to simulate user input. The content
// of the file is then read and compared to the
// expected values
func TestProcessFile(t *testing.T) {
	// Test to read a file
	t.Run("ReadFile", func(t *testing.T) {
		content := "Line 1\nLine 2\nLine 3"
		tempFile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		_, err = tempFile.WriteString(content)
		if err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}

		input, err := cli.ProcessFile(tempFile.Name())
		if err != nil {
			t.Errorf("fxpected no error, but got %v", err)
		}
		if input != content {
			t.Errorf("expected content:\n'%s'\n\nbut got:\n'%s'", content, input)
		}
	})

	// Test to read a file that does not exist
	t.Run("FileNotFound", func(t *testing.T) {
		_, err := cli.ProcessFile("nonexistent.txt")
		if err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	// Test to read an empty file
	t.Run("EmptyFile", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		input, err := cli.ProcessFile(tempFile.Name())
		if err != nil {
			t.Errorf("expected no error, but got: %v", err)
		}
		if input != "" {
			t.Errorf("expected empty input, but got: %s", input)
		}
	})
}
