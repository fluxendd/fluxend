package commands

import (
	"fluxend/internal/app"
	"fluxend/internal/database/seeders"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// seedCmd represents the command to seed the database
var seedCmd = &cobra.Command{
	Use:   "seed [seeder1,seeder2,...]",
	Short: "Seed the database with initial data",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		var customSeeders map[string]bool
		if len(args) > 0 {
			customSeeders = make(map[string]bool)
			crumbs := strings.Split(args[0], ",")
			for _, seeder := range crumbs {
				customSeeders[strings.TrimSpace(seeder)] = true
			}
		}

		runSeeders(customSeeders)
	},
}

func runSeeders(customSeeders map[string]bool) {
	container := app.InitializeContainer()
	log.Info().Msg("Database seeding started")

	seedersToRun := map[string]func(*do.Injector){
		"setting": seeders.Settings,
		"user":    seeders.Users,
	}

	for name, seeder := range seedersToRun {
		if len(customSeeders) > 0 && !customSeeders[name] {
			log.Info().Str("seeder", name).Msg("Skipping seeder as it is not specified in the arguments")
			continue
		}

		log.Info().Str("seeder", name).Msg("Running seeder")
		seeder(container)
	}

	log.Info().Msg("Database seeding completed")
}
