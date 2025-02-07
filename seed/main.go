package main

import (
	"fluxton/di"
	"fluxton/seeders"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"time"
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

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Status    string    `db:"status"`
	RoleID    int       `json:"role_id"`
	Bio       string    `json:"bio"`
	Password  string    `json:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
