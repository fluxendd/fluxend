package cmd

import (
	"fluxton/middlewares"
	"fluxton/repositories"
	"fluxton/routes"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Fluxton API server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	e := setupServer(InitializeContainer())
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func setupServer(container *do.Injector) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e.Use(sentryecho.New(sentryecho.Options{}))

	registerRoutes(e, container)

	return e
}

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
	routes.RegisterFunctionRoutes(e, container, authMiddleware)
	routes.RegisterBackup(e, container, authMiddleware)

	e.GET("/docs/*", echoSwagger.WrapHandler)
}
