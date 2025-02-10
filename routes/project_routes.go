package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterProjectRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, ProjectController *controllers.ProjectController) {
	projectsGroup := e.Group("api/projects", authMiddleware)

	projectsGroup.POST("", ProjectController.Store)
	projectsGroup.GET("", ProjectController.List)
	projectsGroup.GET("/:projectID", ProjectController.Show)
	projectsGroup.PUT("/:projectID", ProjectController.Update)
	projectsGroup.DELETE("/:projectID", ProjectController.Delete)
}
