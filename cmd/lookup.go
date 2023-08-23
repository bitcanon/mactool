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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/mac"
	"github.com/bitcanon/mactool/oui"
)

// lookupAction extracts MAC addresses from the input string,
// performs vendor lookup, and prints the result to the output writer.
func lookupAction(out io.Writer, csv io.Reader, s string) error {
	// Extract MAC addresses from string
	macs, err := mac.FindAllMacAddresses(s)
	if err != nil {
		return err
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
			fmt.Fprintln(out, macAddress)
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

While operating in interactive mode, enter or paste the input string and then press
Enter to proceed. To exit, use Ctrl+D (Unix) or Ctrl+Z (Windows), followed by Enter.`

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

				// Create the CSV file for writing
				outfile, err := os.Create(csv)
				if err != nil {
					return err
				}

				// Download the OUI database
				if err := oui.DownloadDatabase(outfile, url); err != nil {
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

		// Extract MAC addresses from string and
		// perform vendor lookup on each address
		return lookupAction(os.Stdout, file, input)
	},
}

// init registers the lookup command and flags
func init() {
	// Add the lookup command to the root command
	rootCmd.AddCommand(lookupCmd)

	// Add the --csv-file flag to the lookup command
	lookupCmd.PersistentFlags().StringP("csv-file", "f", "oui.csv", "path to CSV file")

	// Check for environment variables prefixed with MACTOOL
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("MACTOOL")

	// Bind the --csv-file flag to the viper variable
	viper.BindPFlag("csv-file", lookupCmd.PersistentFlags().Lookup("csv-file"))
}
