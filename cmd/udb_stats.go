package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// udbStats runs in background and populates table with stats
var udbStats = &cobra.Command{
	Use:   "udb.stats",
	Short: "Pull stats from the user databases",
	Run: func(cmd *cobra.Command, args []string) {
		getDatabaseStats()
	},
}

// TODO: implement later
func getDatabaseStats() {
	// container := InitializeContainer()

	log.Println("Starting database stat compilation...")
}
