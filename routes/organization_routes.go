package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterOrganizationRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, organizationController *controllers.OrganizationController, organizationUserController *controllers.OrganizationUserController) {
	organizationsGroup := e.Group("api/organizations", authMiddleware)

	organizationsGroup.POST("", organizationController.Store)
	organizationsGroup.GET("", organizationController.List)
	organizationsGroup.GET("/:organizationUUID", organizationController.Show)
	organizationsGroup.PUT("/:organizationUUID", organizationController.Update)
	organizationsGroup.DELETE("/:organizationUUID", organizationController.Delete)

	// organization users
	organizationsGroup.POST("/:organizationUUID/users", organizationUserController.Store)
	organizationsGroup.GET("/:organizationUUID/users", organizationUserController.List)
	organizationsGroup.DELETE("/:organizationUUID/users/:userID", organizationUserController.Delete)
}
