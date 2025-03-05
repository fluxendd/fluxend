package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterOrganizationRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	organizationController := do.MustInvoke[*controllers.OrganizationController](container)
	organizationMemberController := do.MustInvoke[*controllers.OrganizationMemberController](container)

	organizationsGroup := e.Group("api/organizations", authMiddleware)

	organizationsGroup.POST("", organizationController.Store)
	organizationsGroup.GET("", organizationController.List)
	organizationsGroup.GET("/:organizationUUID", organizationController.Show)
	organizationsGroup.PUT("/:organizationUUID", organizationController.Update)
	organizationsGroup.DELETE("/:organizationUUID", organizationController.Delete)

	// organization members
	organizationsGroup.POST("/:organizationUUID/members", organizationMemberController.Store)
	organizationsGroup.GET("/:organizationUUID/members", organizationMemberController.List)
	organizationsGroup.DELETE("/:organizationUUID/members/:userID", organizationMemberController.Delete)
}
