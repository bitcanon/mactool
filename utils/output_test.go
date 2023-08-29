package utils_test

import (
	"os"
	"testing"

	"github.com/bitcanon/mactool/utils"
)

// TestGetOutputStream tests the GetOutputStream function
func TestGetOutputStream(t *testing.T) {
	// Write to stdout
	t.Run("WriteToStdout", func(t *testing.T) {
		outStream, err := utils.GetOutputStream("", false)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if outStream != os.Stdout {
			t.Errorf("Expected stdout, but got: %v", outStream)
		}
	})

	// Write to file, overwrite
	t.Run("CreateNewFile", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		outStream, err := utils.GetOutputStream(tempFile.Name(), false)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if outStream == nil {
			t.Errorf("Expected non-nil output stream")
		}
		defer outStream.Close()
	})

	// Write to file, append
	t.Run("AppendToFile", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		outStream, err := utils.GetOutputStream(tempFile.Name(), true)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if outStream == nil {
			t.Errorf("Expected non-nil output stream")
		}
		defer outStream.Close()
	})

	// Invalid filename
	t.Run("InvalidFilename", func(t *testing.T) {
		outStream, err := utils.GetOutputStream("/invalid/path", false)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
		if outStream != nil {
			t.Errorf("Expected nil output stream")
		}
	})
}
