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
	"strings"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/mac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createMacFormatFromFlags(upper bool, lower bool, delimiter string, groupSize int) mac.MacFormat {
	// Select character case based on flags
	caseOption := mac.OriginalCase
	if upper {
		caseOption = mac.Upper
	} else if lower {
		caseOption = mac.Lower
	}

	// Select delimiter based on flags
	delimiterOption := mac.OriginalDelim
	switch delimiter {
	case ":":
		delimiterOption = mac.Colon
	case "-":
		delimiterOption = mac.Hyphen
	case ".":
		delimiterOption = mac.Dot
	case "":
		delimiterOption = mac.None
	}

	// Create a MacFormat struct from the flags
	return mac.MacFormat{
		Case:      caseOption,
		Delimiter: delimiterOption,
		GroupSize: groupSize,
	}
}

// formatAction extracts MAC addresses from the input string,
// performs reformatting according to the provided format,
// and prints the result to the output writer.
func formatAction(out io.Writer, format mac.MacFormat, s string) error {
	// Extract MAC addresses from string
	macs, err := mac.FindAllMacAddresses(s)
	if err != nil {
		return err
	}

	// Print MAC addresses found in the input string
	// to the output writer
	for _, macAddress := range macs {
		// Format the MAC address
		formattedMacAddress, err := mac.FormatMacAddress(macAddress, format)
		if err != nil {
			return err
		}

		// Print the formatted MAC address
		fmt.Fprintln(out, formattedMacAddress)
	}

	// No errors occurred
	return nil
}

// Example help text for the format command
const formatExample = `  mactool format 00:00:5e:00:53:01 --lower-case --delimiter . --group-size 4
  mactool format First address 0000.5E00.5301, second address 00:00:5e:00:53:01, etc. -u -d - -g 2
  cat macs.txt | mactool format --lower-case --delimiter :
  ip addr | mactool format

Interactive mode:
  mactool format

Use interactive mode when you intend to conveniently paste and
format multiple MAC addresses from external sources.`

// Long help text for the format command
const formatLong = `The format command extracts MAC addresses from the input string,
performs reformatting according to the provided argument flags,
and prints the result to the terminal.

The command takes input in the form of command line arguments,
standard input (piped data) or interactive input.`

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:          "format [input]",
	Short:        "Change format of MAC addresses from the input string",
	Long:         formatLong,
	Example:      formatExample,
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

		// Create a MacFormat struct from the flags
		format := createMacFormatFromFlags(
			viper.GetBool("upper-case"),
			viper.GetBool("lower-case"),
			viper.GetString("delimiter"),
			viper.GetInt("group-size"),
		)

		// Format the MAC addresses found in the input string
		// using the format specified by the flags
		return formatAction(os.Stdout, format, input)
	},
}

func init() {
	// Add the format command to the root command
	rootCmd.AddCommand(formatCmd)

	// Persistent flags
	formatCmd.Flags().BoolP("upper-case", "u", false, "convert MAC addresses to upper case")
	formatCmd.Flags().BoolP("lower-case", "l", false, "convert MAC addresses to lower case")
	formatCmd.Flags().StringP("delimiter", "d", ":", "delimiter character to use between hex groups")
	formatCmd.Flags().IntP("group-size", "g", 2, "number of characters in each hex group")

	// Check for environment variables prefixed with MACTOOL
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("MACTOOL")

	// Bind environment variables to flags
	viper.BindPFlag("upper-case", formatCmd.Flags().Lookup("upper-case"))
	viper.BindPFlag("lower-case", formatCmd.Flags().Lookup("lower-case"))
	viper.BindPFlag("delimiter", formatCmd.Flags().Lookup("delimiter"))
	viper.BindPFlag("group-size", formatCmd.Flags().Lookup("group-size"))
}
