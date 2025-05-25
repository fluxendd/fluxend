package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterProjectRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc, allowProjectMiddleware echo.MiddlewareFunc) {
	projectController := do.MustInvoke[*handlers.ProjectHandler](container)

	projectsGroup := e.Group("api/projects", authMiddleware, allowProjectMiddleware)

	projectsGroup.POST("", projectController.Store)
	projectsGroup.GET("", projectController.List)
	projectsGroup.GET("/:projectUUID", projectController.Show)
	projectsGroup.PUT("/:projectUUID", projectController.Update)
	projectsGroup.DELETE("/:projectUUID", projectController.Delete)
}
