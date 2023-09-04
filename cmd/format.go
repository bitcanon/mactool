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
	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createMacFormatFromFlags creates a MacFormat struct from the flags.
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

	// Select group size based on flags
	groupSizeOption := mac.OriginalGroupSize
	switch groupSize {
	case 2:
		groupSizeOption = mac.GroupSizeTwo
	case 4:
		groupSizeOption = mac.GroupSizeFour
	case 6:
		groupSizeOption = mac.GroupSizeSix
	}

	// Create a MacFormat struct from the flags
	return mac.MacFormat{
		Case:      caseOption,
		Delimiter: delimiterOption,
		GroupSize: groupSizeOption,
	}
}

// formatAction finds and formats MAC addresses in the input string.
// The MAC addresses are formatted according to the provided format,
// inside the input string, and printed to the output writer.
func formatAction(out io.Writer, format mac.MacFormat, s string) error {
	// If the input string is empty, return without doing anything
	if len(s) == 0 {
		return nil
	}

	// Split the input string into lines
	lines := strings.Split(s, "\n")

	// Process each line separately
	for i, line := range lines {
		// Find all MAC addresses in the line
		macs, err := mac.FindAllMacAddresses(line)
		if err != nil {
			return err
		}

		// Loop through each MAC address found in the line
		for _, m := range macs {
			// Format the MAC address
			formattedMacAddress, err := mac.FormatMacAddress(m, format)
			if err != nil {
				return err
			}
			// Replace the MAC address with the formatted version
			lines[i] = strings.ReplaceAll(lines[i], m, formattedMacAddress)
		}

		// Print the line to the output writer
		fmt.Fprintln(out, lines[i])
	}

	// No errors occurred
	return nil
}

// Example help text for the format command
const formatExample = `  mactool format 00:00:5e:00:53:01 --lower --delimiter . --group-size 4
  mactool format First address 0000.5E00.5301, second address 00:00:5e:00:53:01, etc. -u -d - -g 2
  cat macs.txt | mactool format --lower --delimiter :
  ip addr | mactool format

Interactive mode:
  mactool format

Use interactive mode when you intend to conveniently paste and
process output from a network device containing MAC addresses.`

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

		// Check if data is being piped, read from file or redirected to stdin
		if viper.GetString("format.input-file") != "" {
			// Read input from file
			input, err = cli.ProcessFile(viper.GetString("format.input-file"))
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

		// Create a MacFormat struct from the flags
		format := createMacFormatFromFlags(
			viper.GetBool("format.upper"),
			viper.GetBool("format.lower"),
			viper.GetString("format.delimiter"),
			viper.GetInt("format.group-size"),
		)

		// Determine the output file using Viper
		outputFile := viper.GetString("format.output-file")
		append := viper.GetBool("format.append")

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

		// Format the MAC addresses found in the input string
		// using the format specified by the flags
		return formatAction(outStream, format, input)
	},
}

func init() {
	// Add the format command to the root command
	rootCmd.AddCommand(formatCmd)

	// Add the --upper flag to the format command
	formatCmd.Flags().BoolP("upper", "u", false, "convert MAC addresses to upper case")
	viper.BindPFlag("format.upper", formatCmd.Flags().Lookup("upper"))

	// Add the --lower flag to the format command
	formatCmd.Flags().BoolP("lower", "l", false, "convert MAC addresses to lower case")
	viper.BindPFlag("format.lower", formatCmd.Flags().Lookup("lower"))

	// Add the --delimiter flag to the format command
	formatCmd.Flags().StringP("delimiter", "d", "=", "delimiter character to use between hex groups")
	viper.BindPFlag("format.delimiter", formatCmd.Flags().Lookup("delimiter"))

	// Add the --group-size flag to the format command
	formatCmd.Flags().IntP("group-size", "g", 0, "number of characters in each hex group")
	viper.BindPFlag("format.group-size", formatCmd.Flags().Lookup("group-size"))

	// Add flag for input file path
	formatCmd.Flags().StringP("input-file", "i", "", "read input from file")
	viper.BindPFlag("format.input-file", formatCmd.Flags().Lookup("input-file"))

	// Add flag for output file path
	formatCmd.Flags().StringP("output-file", "o", "", "write output to file")
	viper.BindPFlag("format.output-file", formatCmd.Flags().Lookup("output-file"))

	// Set to the value of the --append flag if set
	formatCmd.Flags().BoolP("append", "a", false, "append when writing to file with --output-file")
	viper.BindPFlag("format.append", formatCmd.Flags().Lookup("append"))

	// Check for environment variables prefixed with MACTOOL
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("MACTOOL")
}
