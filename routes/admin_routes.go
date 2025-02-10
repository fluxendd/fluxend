package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterAdminRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, settingController *controllers.SettingController) {
	adminGroup := e.Group("api/admin", authMiddleware)

	// settings
	adminGroup.GET("/settings", settingController.List)
	adminGroup.PUT("/settings", settingController.Update)
	adminGroup.PUT("/settings/reset", settingController.Reset)
}
