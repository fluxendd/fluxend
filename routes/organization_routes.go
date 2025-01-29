package routes

import (
	"github.com/labstack/echo/v4"
	"myapp/controllers"
	"myapp/middleware"
)

func RegisterOrganizationRoutes(e *echo.Echo, organizationController *controllers.OrganizationController) {
	notesGroup := e.Group("api/organizations", middleware.AuthMiddleware)

	notesGroup.POST("", organizationController.Store)
	notesGroup.GET("", organizationController.List)
	notesGroup.GET("/:id", organizationController.Show)
	notesGroup.PUT("/:id", organizationController.Update)
	notesGroup.DELETE("/:id", organizationController.Delete)
}
