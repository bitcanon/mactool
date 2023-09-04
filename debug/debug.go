package debug

import (
	"fmt"
	"os"

	"github.com/bitcanon/mactool/oui"
	"github.com/bitcanon/mactool/utils"
	"github.com/spf13/viper"
)

// PrintConfigDebug prints full debug information about the configuration file
// and the variables set in the environment
func PrintConfigDebug() {
	// Get and print the default config file path
	utils.PrintConfigInfo()
	fmt.Println()

	// Print all configuration variables
	utils.PrintVariables(os.Stdout, utils.All)
}

// PrintDatabaseDebug prints debug information about the OUI database file
func PrintDatabaseDebug() {
	// Get the path to the OUI database file
	dbPath := viper.GetString("lookup.oui-file")

	// Get the download URL for the OUI database file
	dbURL := viper.GetString("lookup.oui-url")

	// Get the number of days since the database file was last modified
	days, err := utils.DaysSinceLastModified(dbPath)
	if err != nil {
		fmt.Printf("Failed to get last modified time for %s: %v\n", dbPath, err)
		return
	}

	// Open the database file
	file, err := os.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", dbPath, err)
		return
	}
	defer file.Close()

	// Load the database file
	db, err := oui.LoadDatabase(file)

	// Get the number of entries in the database
	entries := db.GetNumberOfEntries()

	// Print the database file path and the number of days since it was last modified
	fmt.Println("OUI Database:")
	fmt.Printf(" CSV database file URL    : %s\n", dbURL)
	fmt.Printf(" CSV database file path   : %s\n", dbPath)
	fmt.Printf(" Days since last modified : %d\n", days)
	fmt.Printf(" Number of entries        : %d\n", entries)
}
