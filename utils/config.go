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

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

// VariableScope defines the scope of the variables to be returned
type VariableScope int

const (
	All VariableScope = iota
	Environment
	ConfigFile
)

// getLongestKeyLength returns the length of the longest key in a map.
// This is used to align the values when printing the variables.
func getLongestKeyLength(m map[string]string) int {
	maxKeyLength := 0
	for key := range m {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}
	return maxKeyLength
}

// getEnvironmentVariables returns a map containing all the environment variables
// avaliable to the process and filtering out the ones that don't start with the
// specified prefix "MACTOOL_".
func getEnvironmentVariables() map[string]string {
	// We only want to load environment variables that start with this prefix
	prefix := "MACTOOL_"

	// Get all environment variables
	envVars := os.Environ()
	sort.Strings(envVars)

	// Create a map to store the filtered environment variables
	filteredVars := make(map[string]string)

	// Loop through all environment variables
	for _, envVar := range envVars {
		// Split the environment variable into key and value.
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			// Check if the key starts with the specified prefix.
			if strings.HasPrefix(key, prefix) {
				filteredVars[key] = value
			}
		}
	}

	// Return the filtered environment variables
	return filteredVars
}

// GetVariables returns a map containing all the configuration variables
// that are defined in the config file and/or environment variables.
func GetVariables(scope VariableScope) (map[string]string, int) {
	// Define a map to store the configuration variables
	vars := make(map[string]string)

	// If the scope is Environment, return the environment variables
	if scope == Environment {
		vars = getEnvironmentVariables()
		maxKeyLength := getLongestKeyLength(vars)
		return vars, maxKeyLength
	}

	// If the scope is ConfigFile or All, get the configuration variables
	keys := viper.AllKeys()
	sort.Strings(keys)

	// Loop through all configuration variables
	for _, key := range keys {
		value := viper.Get(key)

		// Check if the key is defined in the config file
		// and add it to the map if it is.
		if scope == ConfigFile && viper.InConfig(key) {
			vars[key] = fmt.Sprintf("%v", value)
		} else if scope == All {
			vars[key] = fmt.Sprintf("%v", value)
		}
	}

	// Get the length of the longest key in the map
	maxKeyLength := getLongestKeyLength(vars)

	// Return the filtered configuration variables and the length of the longest key
	return vars, maxKeyLength
}

// GetMapKeys returns a slice containing all the keys from a map.
func GetMapKeys(m map[string]string) []string {
	// Create a slice to store the keys
	keys := make([]string, 0, len(m))

	// Loop through all keys in the map and append them to the slice
	for key := range m {
		keys = append(keys, key)
	}

	// Return the slice of keys
	return keys
}

// PrintVariables prints the variables that are defined in the config file
// and/or environment variables. The scope argument defines which variables
// to print (Environment, ConfigFile or All).
func PrintVariables(out io.Writer, scope VariableScope) error {
	// Get all environment variables.
	vars, maxKeyLength := make(map[string]string), 0

	// Get the variables based on the scope
	srcString := ""
	switch scope {
	case Environment:
		vars, maxKeyLength = GetVariables(Environment)
		srcString = "Environment"
	case ConfigFile:
		vars, maxKeyLength = GetVariables(ConfigFile)
		srcString = "Configuration file"
	case All:
		vars, maxKeyLength = GetVariables(All)
		srcString = "All"
	default:
		return fmt.Errorf("invalid variable scope: %v", scope)
	}

	// Get the keys and sort them so they are printed in alphabetical order
	keys := GetMapKeys(vars)
	sort.Strings(keys)

	// Print the filtered environment variables
	fmt.Fprintln(out, srcString, "variables loaded:")
	for _, key := range keys {
		padding := strings.Repeat(" ", maxKeyLength-len(key))
		fmt.Fprintf(out, " %s%s : %v\n", key, padding, vars[key])
	}

	// Print a message if no environment variables are defined in the config file
	if len(vars) == 0 {
		fmt.Fprintln(out, " No variables defined.")
	}

	// No errors occurred
	return nil
}

// PrintConfigInfo prints the configuration file path, the variables that are
// defined in the config file and the environment variables that start with
// the specified prefix "MACTOOL_".
func PrintConfigInfo() {
	// Get and print the default config file path
	configFilePath := viper.ConfigFileUsed()
	fmt.Printf("Configuration file path: \n %s\n", configFilePath)
	fmt.Println()

	// Print the variables that are defined in the config file
	PrintVariables(os.Stdout, ConfigFile)
	fmt.Println()

	// Print the environment variables that start with the specified prefix
	PrintVariables(os.Stdout, Environment)
}
