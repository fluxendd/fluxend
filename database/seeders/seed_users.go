package seeders

import (
	"fluxton/database/factories"
	"fluxton/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"os"
)

func SeedUsers(container *do.Injector) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	userFactory := do.MustInvoke[*factories.UserFactory](container)

	_, err := userFactory.Create(
		userFactory.WithUsername("superman"),
		userFactory.WithRole(models.UserRoleSuperman),
		userFactory.WithEmail("superman@fluxton.com"),
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
