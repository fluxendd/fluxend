package main

import (
	"flag"
	"fluxton/controllers"
	"fluxton/di"
	"fluxton/routes"
	"fluxton/seeders"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
)

func main() {
	// Define the flag for choosing between "server" and "seed"
	mode := flag.String("cmd", "server", "Possible values: 'server', 'seed'")
	flag.Parse()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the DI container
	container := di.InitializeContainer()

	if *mode == "seed" {
		// Run the seeders if the mode is "seed"
		runSeeders(container)
	} else {
		// Otherwise, start the server
		startServer(container)
	}
}

func startServer(container *do.Injector) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Controllers
	userController := do.MustInvoke[*controllers.UserController](container)
	organizationController := do.MustInvoke[*controllers.OrganizationController](container)
	organizationUserController := do.MustInvoke[*controllers.OrganizationUserController](container)
	projectController := do.MustInvoke[*controllers.ProjectController](container)

	tableController := do.MustInvoke[*controllers.TableController](container)
	columnController := do.MustInvoke[*controllers.ColumnController](container)
	rowController := do.MustInvoke[*controllers.RowController](container)
	indexController := do.MustInvoke[*controllers.IndexController](container)

	// Register routes
	routes.RegisterUserRoutes(e, userController)
	routes.RegisterOrganizationRoutes(e, organizationController, organizationUserController)
	routes.RegisterProjectRoutes(e, projectController)
	routes.RegisterTableRoutes(e, tableController, columnController, indexController)
	routes.RegisterRowRoutes(e, rowController)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func runSeeders(container *do.Injector) {
	log.Info("Starting database seeding...")

	// List of seeders to run
	seedersToRun := []func(*do.Injector){
		seeders.SeedUsers,
	}

	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Info("Database seeding completed successfully.")
}
