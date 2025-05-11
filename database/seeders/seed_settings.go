package seeders

import (
	"fluxton/constants"
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
		{Name: "storageDriver", Value: os.Getenv("STORAGE_DRIVER"), DefaultValue: constants.StorageDriverFilesystem},
		{Name: "maxProjectsPerOrg", Value: "10", DefaultValue: "10"},
		{Name: "allowRegistrations", Value: "yes", DefaultValue: "yes"},
		{Name: "allowProjects", Value: "yes", DefaultValue: "yes"},
		{Name: "allowForms", Value: "yes", DefaultValue: "yes"},
		{Name: "allowStorage", Value: "yes", DefaultValue: "yes"},
		{Name: "allowBackups", Value: "yes", DefaultValue: "yes"},

		// Storage settings
		{Name: "storageMaxContainers", Value: "10", DefaultValue: "10"},
		{Name: "storageMaxFileSizeInKB", Value: "1024", DefaultValue: "1024"},
		{Name: "storageAllowedMimes", Value: "jpg,png,pdf", DefaultValue: "jpg,png,pdf"},

		// API throttle settings
		{Name: "apiThrottleLimit", Value: "100", DefaultValue: "100"},
		{Name: "apiThrottleInterval", Value: "60", DefaultValue: "60"},
		{Name: "allowApiThrottle", Value: "yes", DefaultValue: "no"},

		// External services settings          sss
		{Name: "awsAccessKeyId", Value: os.Getenv("AWS_ACCESS_KEY_ID"), DefaultValue: ""},
		{Name: "awsSecretAccessKey", Value: os.Getenv("AWS_SECRET_ACCESS_KEY"), DefaultValue: ""},
		{Name: "awsRegion", Value: os.Getenv("AWS_REGION"), DefaultValue: "eu-central-1"},
		{Name: "backblazeKeyId", Value: os.Getenv("BACKBLAZE_KEY_ID"), DefaultValue: ""},
		{Name: "backblazeApplicationKey", Value: os.Getenv("BACKBLAZE_APPLICATION_KEY"), DefaultValue: ""},
		{Name: "dropboxAccessToken", Value: os.Getenv("DROPBOX_ACCESS_TOKEN"), DefaultValue: ""},
		{Name: "dropboxAppKey", Value: os.Getenv("DROPBOX_APP_KEY"), DefaultValue: ""},
		{Name: "sendgridApiKey", Value: os.Getenv("SENDGRID_API_KEY"), DefaultValue: ""},
		{Name: "sendgridEmailSource", Value: os.Getenv("SENDGRID_EMAIL_SOURCE"), DefaultValue: ""},
		{Name: "mailgunApiKey", Value: os.Getenv("MAILGUN_API_KEY"), DefaultValue: ""},
		{Name: "mailgunEmailSource", Value: os.Getenv("MAILGUN_EMAIL_SOURCE"), DefaultValue: ""},
		{Name: "mailgunDomain", Value: os.Getenv("MAILGUN_DOMAIN"), DefaultValue: ""},
		{Name: "mailgunRegion", Value: os.Getenv("MAILGUN_REGION"), DefaultValue: "us"},
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
