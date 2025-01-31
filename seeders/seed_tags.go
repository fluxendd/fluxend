package seeders

import (
	"fluxton/factories"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
)

func SeedTags(container *do.Injector) {
	noteFactory := do.MustInvoke[*factories.TagFactory](container)

	_, err := noteFactory.CreateWithName("Default")
	if err != nil {
		log.Fatalf("Error seeding tags: %v", err)
	}

	log.Info("Tags seeded successfully")
}
