package routes

import (
	"fluxton/controllers"
	"fluxton/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterAdminRoutes(e *echo.Echo, settingController *controllers.SettingController) {
	adminGroup := e.Group("api/admin", middlewares.AuthMiddleware)

	// settings
	adminGroup.GET("/settings", settingController.List)
	adminGroup.PUT("/settings", settingController.Update)
	adminGroup.PUT("/settings/reset", settingController.Reset)
}
