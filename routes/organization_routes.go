package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterOrganizationRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, organizationController *controllers.OrganizationController, organizationUserController *controllers.OrganizationUserController) {
	organizationsGroup := e.Group("api/organizations", authMiddleware)

	organizationsGroup.POST("", organizationController.Store)
	organizationsGroup.GET("", organizationController.List)
	organizationsGroup.GET("/:organizationID", organizationController.Show)
	organizationsGroup.PUT("/:organizationID", organizationController.Update)
	organizationsGroup.DELETE("/:organizationID", organizationController.Delete)

	// organization users
	organizationsGroup.POST("/:organizationID/users", organizationUserController.Store)
	organizationsGroup.GET("/:organizationID/users", organizationUserController.List)
	organizationsGroup.DELETE("/:organizationID/users/:userID", organizationUserController.Delete)
}
