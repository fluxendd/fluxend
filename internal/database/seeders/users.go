package seeders

import (
	"fluxend/internal/config/constants"
	"fluxend/internal/database/factories"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
)

func Users(container *do.Injector) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	userFactory := do.MustInvoke[*factories.UserFactory](container)

	_, err := userFactory.Create(
		userFactory.WithUsername(os.Getenv("SUPERUSER_USERNAME")),
		userFactory.WithEmail(os.Getenv("SUPERUSER_EMAIL")),
		userFactory.WithPassword(os.Getenv("SUPERUSER_PASSWORD")),
		userFactory.WithRole(constants.UserRoleSuperman),
	)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Error creating superman user")
	}

	_, err = userFactory.CreateMany(3)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Error creating users")
	}

	log.Info().
		Msg("Users seeded successfully")
}
