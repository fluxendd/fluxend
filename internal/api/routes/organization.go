package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterOrganizationRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	organizationController := do.MustInvoke[*handlers.OrganizationHandler](container)
	organizationMemberController := do.MustInvoke[*handlers.OrganizationMemberHandler](container)

	organizationsGroup := e.Group("organizations", authMiddleware)

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
