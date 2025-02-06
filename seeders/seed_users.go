package seeders

import (
	"fluxton/factories"
	"fluxton/models"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
)

func SeedUsers(container *do.Injector) {
	userFactory := do.MustInvoke[*factories.UserFactory](container)

	_, err := userFactory.Create(
		userFactory.WithUsername("admin"),
		userFactory.WithRole(models.UserRoleOwner),
		userFactory.WithEmail("admin@fluxton.com"),
	)
	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}

	_, err = userFactory.CreateMany(3)
	if err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	log.Info("Users seeded successfully")
}
