package commands

import (
	"fluxend/internal/api/middlewares"
	"fluxend/internal/api/routes"
	"fluxend/internal/app"
	"fluxend/internal/domain/logging"
	"fluxend/internal/domain/setting"
	"fluxend/internal/domain/user"
	"fmt"
	"github.com/getsentry/sentry-go"
	echoSentry "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"strings"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Fluxend API server",
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			if strings.HasSuffix(origin, os.Getenv("APP_URL")) || strings.HasSuffix(origin, os.Getenv("BASE_URL")) {
				return true, nil
			}

			return false, nil
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	if os.Getenv("SENTRY_DSN") != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              os.Getenv("SENTRY_DSN"),
			TracesSampleRate: 1.0,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}

		e.Use(echoSentry.New(echoSentry.Options{}))
	}

	registerRoutes(e, container)

	return e
}

func registerRoutes(e *echo.Echo, container *do.Injector) {
	settingService := do.MustInvoke[setting.Service](container)
	userRepo := do.MustInvoke[user.Repository](container)

	authMiddleware := middlewares.Authentication(userRepo)
	allowProjectMiddleware := middlewares.AllowProject(settingService)
	allowFormMiddleware := middlewares.AllowForm(settingService)
	allowStorageMiddleware := middlewares.AllowStorage(settingService)
	allowBackupMiddleware := middlewares.AllowBackup(settingService)

	requestLogRepo := do.MustInvoke[logging.Repository](container)
	requestLogMiddleware := middlewares.RequestLogger(requestLogRepo)
	e.Use(requestLogMiddleware)

	routes.RegisterUserRoutes(e, container, authMiddleware)
	routes.RegisterAdminRoutes(e, container, authMiddleware)
	routes.RegisterOrganizationRoutes(e, container, authMiddleware)
	routes.RegisterProjectRoutes(e, container, authMiddleware, allowProjectMiddleware)
	routes.RegisterTableRoutes(e, container, authMiddleware)
	routes.RegisterFormRoutes(e, container, authMiddleware, allowFormMiddleware)
	routes.RegisterStorageRoutes(e, container, authMiddleware, allowStorageMiddleware)
	routes.RegisterFunctionRoutes(e, container, authMiddleware)
	routes.RegisterBackup(e, container, authMiddleware, allowBackupMiddleware)

	e.GET("/docs/*", echoSwagger.WrapHandler)
}
