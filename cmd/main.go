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

	e.Logger.Fatal(e.Start(":8080"))
}
