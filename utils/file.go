package utils

import (
	"io"
	"os"
	"time"
)

// DaysSinceLastModified returns the number of days since the specified file
// was last modified. If the file does not exist, an error is returned.
func DaysSinceLastModified(filename string) (int, error) {
	// Get the file info
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	// Get the current time
	currentTime := time.Now()

	// Calculate the duration since the file was last modified
	duration := currentTime.Sub(fileInfo.ModTime())

	// Calculate the number of days
	days := int(duration.Hours() / 24)

	return days, nil
}

// copyFile copies a file from src to dest
func CopyFile(src, dest string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// No errors occurred
	return nil
}
