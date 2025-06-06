package commands

import (
	"fluxend/internal/app"
	"fluxend/internal/database/seeders"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"os"
)

// seedCmd represents the command to seed the database
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		runSeeders()
	},
}

func runSeeders() {
	container := app.InitializeContainer()
	log.Info().Msg("Database seeding started")

	seedersToRun := []func(*do.Injector){
		seeders.Settings,
		seeders.Users,
	}

	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Info().Msg("Database seeding completed")
}
