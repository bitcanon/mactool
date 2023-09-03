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
	"sort"
	"strings"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/mac"
	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// extractAction extracts MAC addresses from the input string
func extractAction(out io.Writer, s string) error {
	// Extract MAC addresses from string
	macs, err := mac.FindAllMacAddresses(s)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Sort MAC addresses in ascending or descending order
	if viper.GetBool("extract.sort-asc") {
		sort.Strings(macs)
	} else if viper.GetBool("extract.sort-desc") {
		sort.Sort(sort.Reverse(sort.StringSlice(macs)))
	}

	// Print MAC addresses found in the input string
	// to the output writer
	for _, mac := range macs {
		fmt.Fprintln(out, mac)
	}

	// No errors occurred
	return nil
}

// Example help text for the extract command
const extractExample = `  mactool extract 0000.5e00.5301 00:00:5e:00:53:01 0000-5e00-5301 00-00-5e-00-53-01
  mactool extract First address 0000.5E00.5301, second address 00:00:5e:00:53:01, etc.
  cat macs.txt | mactool extract
  ipconfig /all | mactool extract

Interactive mode:
  mactool extract

Use interactive mode when you intend to conveniently paste and
process output from a network device containing MAC addresses.`

// Long help text for the extract command
const extractLong = `Extract MAC addresses from the input string

The command takes input in the form of command line arguments,
standard input (piped data) or interactive input.`

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:     "extract [input]",
	Short:   "Extract MAC addresses from the input string",
	Long:    extractLong,
	Example: extractExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Input string to hold the processed input
		var input string
		var err error

		// Check if data is being piped, read from file or redirected to stdin
		if viper.GetString("extract.input-file") != "" {
			// Read input from file
			input, err = cli.ProcessFile(viper.GetString("extract.input-file"))
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

		// Determine the output file using Viper
		outputFile := viper.GetString("extract.output-file")
		append := viper.GetBool("extract.append")

		// Get the output stream
		outStream, err := utils.GetOutputStream(outputFile, append)
		if err != nil {
			return err
		}
		defer outStream.Close()

		// Print the configuration debug if the --debug flag is set
		if viper.GetBool("debug") {
			utils.PrintConfigDebug()
		}

		// Extract MAC addresses from string and
		// print them to standard output
		return extractAction(outStream, input)
	},
}

func init() {
	// Add the extract command to the root command
	rootCmd.AddCommand(extractCmd)

	// Set to the value of the --sort-asc flag if set
	extractCmd.Flags().BoolP("sort-asc", "s", false, "sort output in ascending order")
	viper.BindPFlag("extract.sort-asc", extractCmd.Flags().Lookup("sort-asc"))

	// Set to the value of the --sort-desc flag if set
	extractCmd.Flags().BoolP("sort-desc", "S", false, "sort output in descending order")
	viper.BindPFlag("extract.sort-desc", extractCmd.Flags().Lookup("sort-desc"))

	// Add flag for input file path
	extractCmd.Flags().StringP("input-file", "i", "", "read input from file")
	viper.BindPFlag("extract.input-file", extractCmd.Flags().Lookup("input-file"))

	// Add flag for output file path
	extractCmd.Flags().StringP("output-file", "o", "", "write output to file")
	viper.BindPFlag("extract.output-file", extractCmd.Flags().Lookup("output-file"))

	// Set to the value of the --append flag if set
	extractCmd.Flags().BoolP("append", "a", false, "append when writing to file with --output-file")
	viper.BindPFlag("extract.append", extractCmd.Flags().Lookup("append"))
}
