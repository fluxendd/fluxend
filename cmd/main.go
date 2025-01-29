package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"myapp/controllers"
	"myapp/di"
	"myapp/routes"
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

	// Register routes
	routes.RegisterUserRoutes(e, userController)
	routes.RegisterNoteRoutes(e, noteController)

	e.Logger.Fatal(e.Start(":8080"))
}
