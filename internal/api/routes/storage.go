package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterStorageRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc, allowStorageMiddleware echo.MiddlewareFunc) {
	containerController := do.MustInvoke[*handlers.ContainerHandler](container)
	fileController := do.MustInvoke[*handlers.FileHandler](container)

	projectsGroup := e.Group("api/containers", authMiddleware, allowStorageMiddleware)

	projectsGroup.POST("", containerController.Store)
	projectsGroup.GET("", containerController.List)
	projectsGroup.GET("/:containerUUID", containerController.Show)
	projectsGroup.PUT("/:containerUUID", containerController.Update)
	projectsGroup.DELETE("/:containerUUID", containerController.Delete)

	filesGroup := projectsGroup.Group("/:containerUUID/files")

	filesGroup.POST("", fileController.Store)
	filesGroup.GET("", fileController.List)
	filesGroup.GET("/:fileUUID", fileController.Show)
	filesGroup.PUT("/:fileUUID", fileController.Rename)
	filesGroup.DELETE("/:fileUUID", fileController.Delete)
}
