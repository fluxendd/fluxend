package main

import (
	"flag"
	_ "fluxton/cmd/docs"
	"fluxton/database/seeders"
	"fluxton/middlewares"
	"fluxton/repositories"
	"fluxton/routes"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/do"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"strings"
)

// @title Fluxton API
// @version 1.0
// @description Fluxton is backend as-a-service platform that allows you to build, deploy, and scale applications without managing infrastructure.

// @contact.name API Support
// @contact.url http://github.com/fluxton-io/fluxton
// @contact.email chief@fluxton.io

// @host fluxton.io/api
// @BasePath /v2
func main() {
	mode := flag.String("cmd", "server", "Possible values: 'server', 'seed', 'routes'")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	container := InitializeContainer()
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

	e.Logger.Fatal(e.Start(":80"))
}

// setupServer sets up the Echo server with controllers, routes, and middleware
func setupServer(container *do.Injector) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e.Use(sentryecho.New(sentryecho.Options{}))

	// Register routes
	registerRoutes(e, container)

	return e
}

// registerRoutes registers all routes in the Echo instance
func registerRoutes(e *echo.Echo, container *do.Injector) {
	userRepo := do.MustInvoke[*repositories.UserRepository](container)
	authMiddleware := middlewares.AuthMiddleware(userRepo)

	requestLogRepo := do.MustInvoke[*repositories.RequestLogRepository](container)
	requestLogMiddleware := middlewares.RequestLoggerMiddleware(requestLogRepo)
	e.Use(requestLogMiddleware)

	routes.RegisterUserRoutes(e, container, authMiddleware)
	routes.RegisterAdminRoutes(e, container, authMiddleware)
	routes.RegisterOrganizationRoutes(e, container, authMiddleware)
	routes.RegisterProjectRoutes(e, container, authMiddleware)
	routes.RegisterTableRoutes(e, container, authMiddleware)
	routes.RegisterFormRoutes(e, container, authMiddleware)
	routes.RegisterStorageRoutes(e, container, authMiddleware)

	e.GET("/docs/*", echoSwagger.WrapHandler)
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
