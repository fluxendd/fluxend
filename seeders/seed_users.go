package seeders

import (
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"myapp/factories"
)

func SeedUsers(container *do.Injector) {
	userFactory := do.MustInvoke[*factories.UserFactory](container)

	_, err := userFactory.CreateMany(3)
	if err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	log.Info("Users seeded successfully")
}
