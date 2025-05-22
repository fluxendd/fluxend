package routes

import (
	"fluxton/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterAdminRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	settingHandler := do.MustInvoke[*handlers.SettingHandler](container)
	healthHandler := do.MustInvoke[*handlers.HealthHandler](container)
	logHandler := do.MustInvoke[*handlers.LogHandler](container)

	adminGroup := e.Group("api/admin", authMiddleware)

	// settings
	adminGroup.GET("/settings", settingHandler.List)
	adminGroup.PUT("/settings", settingHandler.Update)
	adminGroup.PUT("/settings/reset", settingHandler.Reset)

	// Health check
	adminGroup.GET("/health", healthHandler.Pulse)

	// Logs
	adminGroup.GET("/logs", logHandler.List)
}
