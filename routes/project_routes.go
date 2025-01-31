package routes

import (
	"github.com/labstack/echo/v4"
	"myapp/controllers"
	"myapp/middleware"
)

func RegisterProjectRoutes(e *echo.Echo, ProjectController *controllers.ProjectController) {
	notesGroup := e.Group("api/projects", middleware.AuthMiddleware)

	notesGroup.POST("", ProjectController.Store)
	notesGroup.GET("", ProjectController.List)
	notesGroup.GET("/:id", ProjectController.Show)
	notesGroup.PUT("/:id", ProjectController.Update)
	notesGroup.DELETE("/:id", ProjectController.Delete)
}
