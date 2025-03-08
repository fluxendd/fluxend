package cmd

import (
	"fluxton/services"
	"github.com/spf13/cobra"
	"log"
)

// udbStats runs in background and populates table with stats
var udbStats = &cobra.Command{
	Use:   "udb.stats [project_uuid]",
	Short: "Pull stats for a project database",
	Run: func(cmd *cobra.Command, args []string) {
		getDatabaseStats()
	},
}

// TODO: implement later
func getDatabaseStats() {
	container := InitializeContainer()

	databaseStatsService := services.NewDatabaseStatsService()
}
