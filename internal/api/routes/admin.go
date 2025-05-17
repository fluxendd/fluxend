package routes

import (
	"fluxton/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterAdminRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	settingController := do.MustInvoke[*handlers.SettingHandler](container)
	healthController := do.MustInvoke[*handlers.HealthHandler](container)

	adminGroup := e.Group("api/admin", authMiddleware)

	// settings
	adminGroup.GET("/settings", settingController.List)
	adminGroup.PUT("/settings", settingController.Update)
	adminGroup.PUT("/settings/reset", settingController.Reset)

	// Health check
	adminGroup.GET("/health", healthController.Pulse)
}
