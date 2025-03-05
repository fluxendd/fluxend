package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterAdminRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	settingController := do.MustInvoke[*controllers.SettingController](container)
	healthController := do.MustInvoke[*controllers.HealthController](container)

	adminGroup := e.Group("api/admin", authMiddleware)

	// settings
	adminGroup.GET("/settings", settingController.List)
	adminGroup.PUT("/settings", settingController.Update)
	adminGroup.PUT("/settings/reset", settingController.Reset)

	// Health check
	adminGroup.GET("/health", healthController.Pulse)
}
