package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// userDatabaseStatsCmd runs in background and populates table with stats
var userDatabaseStatsCmd = &cobra.Command{
	Use:   "udb.stats",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		getDatabaseStats()
	},
}

// TODO: implement later
func getDatabaseStats() {
	// container := InitializeContainer()

	log.Println("Starting database stat compilation...")
}
