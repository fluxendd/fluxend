package cmd

import (
	"fluxton/database/seeders"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"log"
)

// seedCmd represents the command to seed the database
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		runSeeders()
	},
}

func runSeeders() {
	container := InitializeContainer()
	log.Println("Starting database seeding...")

	seedersToRun := []func(*do.Injector){
		seeders.SeedUsers,
	}

	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Println("Database seeding completed successfully.")
}
