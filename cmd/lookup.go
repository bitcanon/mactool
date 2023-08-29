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
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/mac"
	"github.com/bitcanon/mactool/oui"
	"github.com/bitcanon/mactool/utils"
)

// lookupAction extracts MAC addresses from the input string,
// performs vendor lookup, and prints the result to the output writer.
func lookupAction(out io.Writer, csv io.Reader, s string) error {
	// Extract MAC addresses from string
	macs, err := mac.FindAllMacAddresses(s)
	if err != nil {
		return err
	}

	// Sort MAC addresses in ascending or descending order
	if viper.GetBool("lookup.sort-asc") {
		sort.Strings(macs)
	} else if viper.GetBool("lookup.sort-desc") {
		sort.Sort(sort.Reverse(sort.StringSlice(macs)))
	}

	// Load the OUI database into memory
	db, err := oui.LoadDatabase(csv)
	if err != nil {
		return err
	}

	// Print MAC addresses found in the input string
	// to the output writer
	for _, macAddress := range macs {
		// Lookup the vendor of the MAC address
		assignment, err := mac.ExtractOuiFromMac(macAddress)
		if err != nil {
			return err
		}

		// Lookup the vendor in the OUI database
		vendor := db.FindOuiByAssignment(assignment)

		if vendor != nil {
			// If the vendor was found, print the vendor name
			fmt.Fprintf(out, "%s (%s)\n", macAddress, vendor.Organization)
		} else {
			// If the vendor was not found, print the MAC address
			// if the --suppress-unmatched flag is not set
			if !viper.GetBool("suppress-unmatched") {
				fmt.Fprintln(out, macAddress)
			}
		}
	}

	// No errors occurred
	return nil
}

// Example help text for the lookup command
const lookupExample = `  mactool lookup 00:00:5e:00:53:01
  mactool lookup 0000.5e00.5301 00:00:5e:00:53:01 0000-5e00-5301 00-00-5e-00-53-01
  mactool lookup First address 0000.5E00.5301, second address 00:00:5e:00:53:01, etc.
  cat macs.txt | mactool lookup
  ip addr | mactool lookup

Interactive mode:
  mactool lookup

Use interactive mode when you intend to conveniently paste and
process output from a network device containing MAC addresses.`

// Long help text for the lookup command
const lookupLong = `Extract MAC addresses from the input string, perform
vendor lookup, and display the result on the terminal.

The command takes input in the form of command line arguments,
standard input (piped data) or interactive input.`

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:          "lookup [input]",
	Short:        "Lookup vendors of MAC addresses from the input string",
	Long:         lookupLong,
	Example:      lookupExample,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Input string to hold the processed input
		var input string
		var err error

		// Check if data is being piped or redirected to stdin
		if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
			// Process data from pipe or redirection (stdin)
			input, err = cli.ProcessStdin()
			if err != nil {
				return err
			}
		} else {
			if len(args) == 0 {
				// If there are no command line arguments,
				// enter interactive mode and read user input
				input, err = cli.ProcessInteractiveInput()
				if err != nil {
					return err
				}
			} else {
				// If there are command line arguments, join them
				// into a single string and use that as user input
				input = strings.Join(args, " ")
			}
		}

		// Get the OUI database file
		csv := viper.GetString("csv-file")

		// Check if the CSV file exists and download it if it doesn't
		_, err = os.Stat(csv)
		if os.IsNotExist(err) {
			fmt.Printf("The file '%s' could not be found.\n", csv)
			fmt.Print("Would you like to download it? (Y/n): ")

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimRight(input, "\r\n")

			if input == "n" || input == "N" {
				// User cancelled the download so exit the program
				fmt.Println("File download cancelled.")
				return nil
			} else {
				// User confirmed the download so proceed
				url := "http://standards-oui.ieee.org/oui/oui.csv"

				// Create a temporary file for writing the downloaded file
				tempFile, err := os.CreateTemp("", "oui.csv")
				if err != nil {
					fmt.Println("Error creating temporary file:", err)
					return err
				}
				defer os.Remove(tempFile.Name()) // Clean up the temporary file when done

				// Download the OUI database
				if err := oui.DownloadDatabase(tempFile, url); err != nil {
					return err
				}

				// Copy the temporary file to the CSV file
				err = copyFile(tempFile.Name(), csv)
				if err != nil {
					return err
				}
			}
		}

		// Open the CSV file
		file, err := os.Open(csv)
		if err != nil {
			return err
		}
		defer file.Close()

		// Determine the output file using Viper
		outputFile := viper.GetString("lookup.output")
		append := viper.GetBool("lookup.append")

		// Get the output stream
		outStream, err := utils.GetOutputStream(outputFile, append)
		if err != nil {
			return err
		}
		defer outStream.Close()

		// Extract MAC addresses from string and
		// perform vendor lookup on each address
		return lookupAction(outStream, file, input)
	},
}

// init registers the lookup command and flags
func init() {
	// Add the lookup command to the root command
	rootCmd.AddCommand(lookupCmd)

	/*
	  The CSV file can be specified using the following methods,
	  in order of precedence (1 is highest, 4 is lowest):
	    1. --csv-file flag
	    2. MACTOOL_CSV_FILE environment variable
	    3. .mactool.yaml config file
	    4. Default CSV file path
	*/

	// Set a default CSV file path
	viper.SetDefault("csv-file", oui.GetDefaultDatabasePath())

	// Set to environment variable MACTOOL_CSV_FILE if set
	err := viper.BindEnv("csv-file")
	cobra.CheckErr(err)

	// Set default path for the flag help text
	var defaultPath string
	if runtime.GOOS == "windows" {
		defaultPath = "%LOCALAPPDATA%\\Mactool\\oui.csv"
	} else {
		defaultPath = "~/.local/share/mactool/oui.csv"
	}

	// Set to the value of the --csv-file flag if set
	lookupCmd.PersistentFlags().StringP("csv-file", "f", "", "path to CSV file (default "+defaultPath+")")
	viper.BindPFlag("csv-file", lookupCmd.PersistentFlags().Lookup("csv-file"))

	// Set to the value of the --suppress-unmatched flag if set
	lookupCmd.PersistentFlags().BoolP("suppress-unmatched", "u", false, "suppress unmatched MAC addresses from output")
	viper.BindPFlag("suppress-unmatched", lookupCmd.PersistentFlags().Lookup("suppress-unmatched"))

	// Set to the value of the --sort-asc flag if set
	lookupCmd.Flags().BoolP("sort-asc", "s", false, "sort output in ascending order")
	viper.BindPFlag("lookup.sort-asc", lookupCmd.Flags().Lookup("sort-asc"))

	// Set to the value of the --sort-desc flag if set
	lookupCmd.Flags().BoolP("sort-desc", "S", false, "sort output in descending order")
	viper.BindPFlag("lookup.sort-desc", lookupCmd.Flags().Lookup("sort-desc"))

	// Add flag for input file path
	lookupCmd.Flags().StringP("input", "i", "", "read input from file")
	viper.BindPFlag("lookup.input", lookupCmd.Flags().Lookup("input"))

	// Add flag for output file path
	lookupCmd.Flags().StringP("output", "o", "", "write output to file")
	viper.BindPFlag("lookup.output", lookupCmd.Flags().Lookup("output"))

	// Set to the value of the --append flag if set
	lookupCmd.Flags().BoolP("append", "a", false, "append when writing to file with --output")
	viper.BindPFlag("lookup.append", lookupCmd.Flags().Lookup("append"))
}

// copyFile copies a file from src to dest
func copyFile(src, dest string) error {
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
