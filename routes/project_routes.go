package routes

import (
	"fluxton/controllers"
	"fluxton/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterProjectRoutes(e *echo.Echo, ProjectController *controllers.ProjectController) {
	projectsGroup := e.Group("api/projects", middlewares.AuthMiddleware)

	projectsGroup.POST("", ProjectController.Store)
	projectsGroup.GET("", ProjectController.List)
	projectsGroup.GET("/:id", ProjectController.Show)
	projectsGroup.PUT("/:id", ProjectController.Update)
	projectsGroup.DELETE("/:id", ProjectController.Delete)
}
