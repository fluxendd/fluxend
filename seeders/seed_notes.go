package seeders

import (
	"fluxton/factories"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
)

func SeedNotes(container *do.Injector) {
	noteFactory := do.MustInvoke[*factories.NoteFactory](container)

	_, err := noteFactory.CreateMany(5, noteFactory.WithUserId(1))
	if err != nil {
		log.Fatalf("Error seeding notes: %v", err)
	}

	log.Info("Notes seeded successfully")
}
