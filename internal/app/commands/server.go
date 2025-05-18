package commands

import (
	middlewares2 "fluxton/internal/api/middlewares"
	routes2 "fluxton/internal/api/routes"
	"fluxton/internal/app"
	"fluxton/internal/database/repositories"
	"fluxton/internal/domain/setting"
	"fluxton/internal/domain/user"
	//"fluxton/repositories"
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
	e := setupServer(app.InitializeContainer())
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func setupServer(container *do.Injector) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if os.Getenv("SENTRY_DSN") != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              os.Getenv("SENTRY_DSN"),
			TracesSampleRate: 1.0,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}

		e.Use(sentryecho.New(sentryecho.Options{}))
	}

	registerRoutes(e, container)

	return e
}

func registerRoutes(e *echo.Echo, container *do.Injector) {
	settingService := do.MustInvoke[setting.Service](container)
	userRepo := do.MustInvoke[user.Repository](container)

	authMiddleware := middlewares2.AuthMiddleware(userRepo)
	allowProjectMiddleware := middlewares2.AllowProjectMiddleware(settingService)
	allowFormMiddleware := middlewares2.AllowFormMiddleware(settingService)
	allowStorageMiddleware := middlewares2.AllowStorageMiddleware(settingService)
	allowBackupMiddleware := middlewares2.AllowBackupMiddleware(settingService)

	requestLogRepo := do.MustInvoke[*repositories.RequestLogRepository](container)
	requestLogMiddleware := middlewares2.RequestLoggerMiddleware(requestLogRepo)
	e.Use(requestLogMiddleware)

	routes2.RegisterUserRoutes(e, container, authMiddleware)
	routes2.RegisterAdminRoutes(e, container, authMiddleware)
	routes2.RegisterOrganizationRoutes(e, container, authMiddleware)
	routes2.RegisterProjectRoutes(e, container, authMiddleware, allowProjectMiddleware)
	routes2.RegisterTableRoutes(e, container, authMiddleware)
	routes2.RegisterFormRoutes(e, container, authMiddleware, allowFormMiddleware)
	routes2.RegisterStorageRoutes(e, container, authMiddleware, allowStorageMiddleware)
	routes2.RegisterFunctionRoutes(e, container, authMiddleware)
	routes2.RegisterBackup(e, container, authMiddleware, allowBackupMiddleware)

	e.GET("/docs/*", echoSwagger.WrapHandler)
}
