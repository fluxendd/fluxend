package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterFormRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, FormController *controllers.FormController) {
	projectsGroup := e.Group("api/projects/:projectUUID/forms", authMiddleware)

	projectsGroup.POST("", FormController.Store)
	projectsGroup.GET("", FormController.List)
	projectsGroup.GET("/:formUUID", FormController.Show)
	projectsGroup.PUT("/:formUUID", FormController.Update)
	projectsGroup.DELETE("/:formUUID", FormController.Delete)
}
