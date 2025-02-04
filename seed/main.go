package main

import (
	"fluxton/di"
	"fluxton/seeders"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
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
	}

	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Info("Database seeding completed successfully.")
}
