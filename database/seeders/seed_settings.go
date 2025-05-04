package seeders

import (
	"fluxton/models"
	"fluxton/repositories"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
)

func SeedSettings(container *do.Injector) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	settingsService := do.MustInvoke[*repositories.SettingRepository](container)

	settings := []models.Setting{
		// General settings
		{Name: "appTitle", Value: os.Getenv("APP_TITLE"), DefaultValue: os.Getenv("APP_TITLE")},
		{Name: "appUrl", Value: os.Getenv("APP_URL"), DefaultValue: os.Getenv("APP_URL")},
		{Name: "jwtSecret", Value: os.Getenv("JWT_SECRET"), DefaultValue: os.Getenv("JWT_SECRET")},
		{Name: "maxProjectsPerOrg", Value: "10", DefaultValue: "10"},
		{Name: "allowRegistrations", Value: "yes", DefaultValue: "yes"},
		{Name: "allowNewProjects", Value: "yes", DefaultValue: "yes"},
		{Name: "enableForms", Value: "yes", DefaultValue: "yes"},
		{Name: "enableStorage", Value: "yes", DefaultValue: "yes"},
		{Name: "enableBackups", Value: "yes", DefaultValue: "yes"},

		// Storage settings
		{Name: "storageMaxBuckets", Value: "10", DefaultValue: "10"},
		{Name: "storageMaxFileSizeInKB", Value: "1024", DefaultValue: "1024"},
		{Name: "storageAllowedMimes", Value: "jpg,png,pdf", DefaultValue: "jpg,png,pdf"},

		// API throttle settings
		{Name: "apiThrottleLimit", Value: "100", DefaultValue: "100"},
		{Name: "apiThrottleInterval", Value: "60", DefaultValue: "60"},
		{Name: "apiThrottleEnabled", Value: "yes", DefaultValue: "no"},
	}

	_, err := settingsService.CreateMany(settings)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Error creating settings")
		return
	}

	log.Info().Msg("Settings seeded successfully")
}
