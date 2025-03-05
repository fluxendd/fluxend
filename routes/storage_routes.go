package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterStorageRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	bucketController := do.MustInvoke[*controllers.BucketController](container)
	fileController := do.MustInvoke[*controllers.FileController](container)

	projectsGroup := e.Group("api/projects/:projectUUID/buckets", authMiddleware)

	projectsGroup.POST("", bucketController.Store)
	projectsGroup.GET("", bucketController.List)
	projectsGroup.GET("/:bucketUUID", bucketController.Show)
	projectsGroup.PUT("/:bucketUUID", bucketController.Update)
	projectsGroup.DELETE("/:bucketUUID", bucketController.Delete)

	filesGroup := projectsGroup.Group("/:bucketUUID/files")

	filesGroup.POST("", fileController.Store)
	filesGroup.GET("", fileController.List)
	filesGroup.GET("/:fileUUID", fileController.Show)
	filesGroup.PUT("/:fileUUID", fileController.Rename)
	filesGroup.DELETE("/:fileUUID", fileController.Delete)
}
