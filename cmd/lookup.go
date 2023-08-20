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

	"github.com/spf13/cobra"

	"github.com/bitcanon/mactool/cli"
	"github.com/bitcanon/mactool/mac"
)

func lookupAction(out io.Writer, s string) error {
	// Extract MAC addresses from string
	macs, err := mac.FindAllMacAddresses(s)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Print MAC addresses found in the input string
	// to the output writer
	for _, mac := range macs {
		fmt.Fprintln(out, mac)
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

// Long help text for the extract command
const lookupLong = `Extract MAC addresses from the input string, perform
vendor lookup, and display the result on the terminal.

The command takes input in the form of command line arguments,
standard input (piped data) or interactive input.`

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:     "lookup",
	Short:   "Lookup vendors of MAC addresses from the input string",
	Long:    lookupLong,
	Example: lookupExample,
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

		// Extract MAC addresses from string and
		// perform vendor lookup on each address
		return lookupAction(os.Stdout, input)
	},
}

func init() {
	rootCmd.AddCommand(lookupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lookupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lookupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
