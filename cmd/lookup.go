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
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/debug"
	"github.com/bitcanon/mactool/mac"
	"github.com/bitcanon/mactool/oui"
	"github.com/bitcanon/mactool/utils"
)

// lookupAction extracts MAC addresses from the input string,
// performs vendor lookup, and prints the result to the output writer.
func lookupAction(out io.Writer, ouiCsvFile io.Reader, s string) error {
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
	db, err := oui.LoadDatabase(ouiCsvFile)
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
			// Write in CSV format if the --csv flag is set
			if viper.GetBool("lookup.csv") {
				csvRow, err := utils.ConvertStringSliceToCSV([]string{macAddress, vendor.Organization, vendor.Address})
				if err != nil {
					return err
				}
				fmt.Fprint(out, csvRow)
			} else {
				// If the vendor was found, print the vendor name
				fmt.Fprintf(out, "%s (%s)\n", macAddress, vendor.Organization)
			}
		} else {
			// If the vendor was not found, print the MAC address
			// if the --suppress-unmatched flag is not set
			if !viper.GetBool("lookup.suppress-unmatched") {
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

		// Check if data is being piped, read from file or redirected to stdin
		if viper.GetString("lookup.input-file") != "" {
			// Read input from file
			input, err = cli.ProcessFile(viper.GetString("lookup.input-file"))
			if err != nil {
				return err
			}
		} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
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
		csv := viper.GetString("lookup.oui-file")

		// Check if the CSV file exists and download it if it doesn't
		oui.UpdateDatabase(csv)

		// Open the CSV file
		file, err := os.Open(csv)
		if err != nil {
			return err
		}
		defer file.Close()

		// Determine the output file using Viper
		outputFile := viper.GetString("lookup.output-file")
		append := viper.GetBool("lookup.append")

		// Get the output stream
		outStream, err := utils.GetOutputStream(outputFile, append)
		if err != nil {
			return err
		}
		defer outStream.Close()

		// Print the configuration debug if the --debug flag is set
		if viper.GetBool("debug") {
			debug.PrintConfigDebug()
		}

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
	    1. --oui-file flag
	    2. MACTOOL_CSV_FILE environment variable
	    3. .mactool.yaml config file
	    4. Default CSV file path
	*/

	// Set a default CSV file path
	viper.SetDefault("lookup.oui-file", oui.GetDefaultDatabasePath())

	// Set a default URL for the OUI CSV file
	viper.SetDefault("lookup.oui-url", "http://standards-oui.ieee.org/oui/oui.csv")

	// Set default path for the flag help text
	var defaultPath string
	if runtime.GOOS == "windows" {
		defaultPath = "%LOCALAPPDATA%\\Mactool\\oui.csv"
	} else {
		defaultPath = "~/.local/share/mactool/oui.csv"
	}

	// Set to the value of the --oui-file flag if set
	lookupCmd.PersistentFlags().StringP("oui-file", "O", "", "path to OUI CSV file (default "+defaultPath+")")
	viper.BindPFlag("lookup.oui-file", lookupCmd.PersistentFlags().Lookup("oui-file"))

	// Set to the value of the --suppress-unmatched flag if set
	lookupCmd.Flags().BoolP("suppress-unmatched", "u", false, "suppress unmatched MAC addresses from output")
	viper.BindPFlag("lookup.suppress-unmatched", lookupCmd.Flags().Lookup("suppress-unmatched"))

	// Set to the value of the --sort-asc flag if set
	lookupCmd.PersistentFlags().BoolP("sort-asc", "s", false, "sort output in ascending order")
	viper.BindPFlag("lookup.sort-asc", lookupCmd.PersistentFlags().Lookup("sort-asc"))

	// Set to the value of the --sort-desc flag if set
	lookupCmd.PersistentFlags().BoolP("sort-desc", "S", false, "sort output in descending order")
	viper.BindPFlag("lookup.sort-desc", lookupCmd.PersistentFlags().Lookup("sort-desc"))

	// Add flag for --input-file path
	lookupCmd.Flags().StringP("input-file", "i", "", "read input from file")
	viper.BindPFlag("lookup.input-file", lookupCmd.Flags().Lookup("input-file"))

	// Add flag for --output-file path
	lookupCmd.PersistentFlags().StringP("output-file", "o", "", "write output to file")
	viper.BindPFlag("lookup.output-file", lookupCmd.PersistentFlags().Lookup("output-file"))

	// Set to the value of the --append flag if set
	lookupCmd.PersistentFlags().BoolP("append", "a", false, "append when writing to file with --output-file")
	viper.BindPFlag("lookup.append", lookupCmd.PersistentFlags().Lookup("append"))

	// Set to the value of the --csv flag if set
	lookupCmd.PersistentFlags().BoolP("csv", "c", false, "write output in CSV format")
	viper.BindPFlag("lookup.csv", lookupCmd.PersistentFlags().Lookup("csv"))
}
