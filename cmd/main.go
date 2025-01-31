package main

import (
	"fluxton/controllers"
	"fluxton/di"
	"fluxton/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// DI container
	container := di.InitializeContainer()

	// controllers
	userController := do.MustInvoke[*controllers.UserController](container)
	noteController := do.MustInvoke[*controllers.NoteController](container)
	organizationController := do.MustInvoke[*controllers.OrganizationController](container)
	projectController := do.MustInvoke[*controllers.ProjectController](container)
	tableController := do.MustInvoke[*controllers.TableController](container)

	// Register routes
	routes.RegisterUserRoutes(e, userController)
	routes.RegisterNoteRoutes(e, noteController)
	routes.RegisterOrganizationRoutes(e, organizationController)
	routes.RegisterProjectRoutes(e, projectController)
	routes.RegisterTableRoutes(e, tableController)

	e.Logger.Fatal(e.Start(":8080"))
}
