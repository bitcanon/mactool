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
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bitcanon/mactool/debug"
	"github.com/bitcanon/mactool/oui"
	"github.com/bitcanon/mactool/utils"
)

// lookupAction extracts MAC addresses from the input string,
// performs vendor lookup, and prints the result to the output writer.
func lookupVendorAction(out io.Writer, ouiCsvFile io.Reader, s string) error {
	// Load the OUI database into memory
	db, err := oui.LoadDatabase(ouiCsvFile)
	if err != nil {
		return err
	}

	// Setup the filter option
	filterOptions := oui.FilterOptions{
		Assignment:   viper.GetBool("lookup-vendor.assignment"),
		Organization: viper.GetBool("lookup-vendor.organization"),
		Address:      viper.GetBool("lookup-vendor.address"),
	}

	// Find all vendors matching the input string
	vendors, err := db.FindAllVendors(s, filterOptions)
	if err != nil {
		return err
	}

	// Sort MAC addresses in ascending or descending order
	if viper.GetBool("lookup.sort-asc") {
		sort.Sort(vendors)
	} else if viper.GetBool("lookup.sort-desc") {
		sort.Sort(sort.Reverse(vendors))
	}

	// Print the vendors found in the input string
	// to the output writer
	for _, vendor := range vendors.Entries {
		// Write in CSV format if the --csv flag is set
		if viper.GetBool("lookup.csv") {
			csvRow, err := utils.ConvertStringSliceToCSV([]string{vendor.Assignment, vendor.Organization, vendor.Address})
			if err != nil {
				return err
			}
			_, err = out.Write([]byte(csvRow))
			if err != nil {
				return err
			}
		} else {
			// If the vendor was found, print the vendor name
			_, err = out.Write([]byte(vendor.Assignment + " " + vendor.Organization + "\n"))
			if err != nil {
				return err
			}
		}
	}

	// No errors occurred
	return nil
}

// Example help text for the lookup command
const lookupVendorExample = `  mactool lookup vendor cisco
  mactool lookup vendor mikrotik
  mactool lookup vendor --assignment 00000C
  mactool lookup vendor --organization "Cisco Systems"
  mactool lookup vendor --address "San Jose"`

// Long help text for the lookup command
const lookupVendorLong = `Find all the OUIs belonging to a vendor or organization.

By default, the lookup vendor command searches in all columns of the CSV database. 
To search in a specific column, use the appropriate flag (e.g. --assignment).

The search is case-insensitive and matches partial strings.
`

// lookupCmd represents the lookup command
var lookupVendorCmd = &cobra.Command{
	Use:          "vendor [input]",
	Short:        "Find all the OUIs belonging to a vendor or organization",
	Long:         lookupVendorLong,
	Example:      lookupVendorExample,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Input string to hold the processed input
		var input string
		var err error

		// Get the search string from the command line arguments
		input = strings.Join(args, " ")

		// If no search string was specified, print the help text
		if input == "" {
			cmd.Help()
			os.Exit(0)
		}

		// Get the OUI database file
		csv := viper.GetString("lookup.oui-file")

		// Check if the CSV file exists and download it if it doesn't
		oui.UpdateDatabase(csv)

		// Open the CSV file
		csvFile, err := os.Open(csv)
		if err != nil {
			return err
		}
		defer csvFile.Close()

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

		// Perform the lookup
		return lookupVendorAction(outStream, csvFile, input)
	},
}

// init registers the lookup vendor command and flags
func init() {
	// Add the lookup vendor command to the lookup command
	lookupCmd.AddCommand(lookupVendorCmd)

	// Add the --assignment flag to the lookup vendor command
	lookupVendorCmd.Flags().Bool("assignment", false, "search in assignment column (e.g. \"A1B2C3\")")
	viper.BindPFlag("lookup-vendor.assignment", lookupVendorCmd.Flags().Lookup("assignment"))

	// Add the --organization flag to the lookup vendor command
	lookupVendorCmd.Flags().Bool("organization", false, "search in organization column")
	viper.BindPFlag("lookup-vendor.organization", lookupVendorCmd.Flags().Lookup("organization"))

	// Add the --address flag to the lookup vendor command
	lookupVendorCmd.Flags().Bool("address", false, "search in address column")
	viper.BindPFlag("lookup-vendor.address", lookupVendorCmd.Flags().Lookup("address"))
}
