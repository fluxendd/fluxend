package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"myapp/di"
	"myapp/seeders"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	container := di.InitializeContainer()

	// Call seeders
	runSeeders(container)
}

func runSeeders(container *do.Injector) {
	log.Info("Starting database seeding...")

	seedersToRun := []func(*do.Injector){
		seeders.SeedUsers,
		seeders.SeedTags,
		seeders.SeedNotes,
	}

	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Info("Database seeding completed successfully.")
}
