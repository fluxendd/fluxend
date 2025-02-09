package main

import (
	"flag"
	"fluxton/controllers"
	"fluxton/di"
	"fluxton/routes"
	"fluxton/seeders"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	"strings"
)

func main() {
	mode := flag.String("cmd", "server", "Possible values: 'server', 'seed', 'routes'")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	container := di.InitializeContainer()
	handleMode(container, *mode)
}

// handleMode determines which functionality to execute based on the "cmd" flag
func handleMode(container *do.Injector, mode string) {
	switch mode {
	case "seed":
		runSeeders(container)
	case "routes":
		PrintRoutes(container)
	default:
		startServer(container)
	}
}

func startServer(container *do.Injector) {
	e := setupServer(container)

	e.Logger.Fatal(e.Start(":8080"))
}

// setupServer sets up the Echo server with controllers, routes, and middleware
func setupServer(container *do.Injector) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	registerRoutes(e, container)

	return e
}

// registerRoutes registers all routes in the Echo instance
func registerRoutes(e *echo.Echo, container *do.Injector) {
	// Controllers
	userController := do.MustInvoke[*controllers.UserController](container)
	settingController := do.MustInvoke[*controllers.SettingController](container)
	organizationController := do.MustInvoke[*controllers.OrganizationController](container)
	organizationUserController := do.MustInvoke[*controllers.OrganizationUserController](container)
	projectController := do.MustInvoke[*controllers.ProjectController](container)
	tableController := do.MustInvoke[*controllers.TableController](container)
	columnController := do.MustInvoke[*controllers.ColumnController](container)
	rowController := do.MustInvoke[*controllers.RowController](container)
	indexController := do.MustInvoke[*controllers.IndexController](container)

	// Register routes
	routes.RegisterUserRoutes(e, userController)
	routes.RegisterAdminRoutes(e, settingController)
	routes.RegisterOrganizationRoutes(e, organizationController, organizationUserController)
	routes.RegisterProjectRoutes(e, projectController)
	routes.RegisterTableRoutes(e, tableController, columnController, indexController)
	routes.RegisterRowRoutes(e, rowController)
}

// runSeeders runs all seeders defined in the seeders package
func runSeeders(container *do.Injector) {
	log.Info("Starting database seeding...")

	// List of seeders to run
	seedersToRun := []func(*do.Injector){
		seeders.SeedUsers,
	}

	// Execute each seeder
	for _, seeder := range seedersToRun {
		seeder(container)
	}

	log.Info("Database seeding completed successfully.")
}

func PrintRoutes(container *do.Injector) {
	e := setupServer(container)

	registeredRoutes := e.Routes()

	// Map to hold groups of routes based on their prefix
	routesGrouped := make(map[string][]string)

	// Iterate through the routes and group them by prefix
	for _, route := range registeredRoutes {
		if route.Method == "echo_route_not_found" {
			continue
		}

		// Extract the prefix (e.g., /api/projects from /api/projects/:projectID/tables)
		routePrefix := getRoutePrefix(route.Path)

		// Group the routes by prefix
		routesGrouped[routePrefix] = append(routesGrouped[routePrefix], fmt.Sprintf("Method: %-6s  Path: %-30s", route.Method, route.Path))
	}

	// Print the grouped routes
	for prefix, matchedRoutes := range routesGrouped {
		fmt.Printf("\nRoutes for %s:\n", prefix)

		for _, route := range matchedRoutes {
			fmt.Println(route)
		}
	}
}

func getRoutePrefix(path string) string {
	// Split the path by '/' and take the first two segments (or fewer if the path is shorter)
	segments := strings.Split(path, "/")

	// Return the prefix by joining the first segments (e.g., /api/projects)
	if len(segments) > 1 {
		return "/" + segments[2]
	}

	return path
}
