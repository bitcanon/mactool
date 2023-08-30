package utils_test

import (
	"testing"

	"github.com/bitcanon/mactool/utils"
)

// TestConvertStringSliceToCSV tests the ConvertStringSliceToCSV function
// using values containing characters that need to be escaped in CSV
func TestConvertStringSliceToCSV(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "EmptyInput",
			input:    []string{},
			expected: "",
		},
		{
			name:     "SingleValue",
			input:    []string{"value1"},
			expected: "value1\n",
		},
		{
			name:     "MultipleValues",
			input:    []string{"value1", "value2", "value3"},
			expected: "value1,value2,value3\n",
		},
		{
			name:     "MultipleValuesWithCommas",
			input:    []string{"value1", "value,2", "value3"},
			expected: "value1,\"value,2\",value3\n",
		},
		{
			name:     "MultipleValuesWithQuotes",
			input:    []string{"value1", "value\"2", "value3"},
			expected: "value1,\"value\"\"2\",value3\n",
		},
		{
			name:     "MultipleValuesWithQuotesAndCommas",
			input:    []string{"value1", "value,\"2", "value3"},
			expected: "value1,\"value,\"\"2\",value3\n",
		},
	}

	// Loop through test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Convert the string slice to CSV
			csvFormatted, err := utils.ConvertStringSliceToCSV(testCase.input)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			// Check if the CSV-formatted string matches the expected string
			if csvFormatted != testCase.expected {
				t.Errorf("expected:\n'%s'\ngot:\n'%s'", testCase.expected, csvFormatted)
			}
		})
	}
}
