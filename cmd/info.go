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

	"github.com/bitcanon/mactool/debug"
	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoAction prints configuration and database information
func infoAction(out io.Writer, s string) error {
	utils.PrintConfigInfo()
	fmt.Println()
	debug.PrintDatabaseDebug()
	return nil
}

// Long help text for the info command
const infoLong = `Display the configuration settings as set in the configuration file,
environment variables, and default values. Additionally, provide details
about the OUI database file, including the count of entries and the duration
in days since the file's last modification.`

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:          "info",
	Short:        "Print configuration and database information",
	Long:         infoLong,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		return infoAction(os.Stdout, "input-text")
	},
}

func init() {
	// Add the info command to the root command
	rootCmd.AddCommand(infoCmd)

	// Check for environment variables prefixed with MACTOOL
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("MACTOOL")
}
