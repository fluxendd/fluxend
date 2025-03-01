package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterStorageRoutes(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	BucketController *controllers.BucketController,
	FileController *controllers.FileController,
) {
	projectsGroup := e.Group("api/projects/:projectUUID/buckets", authMiddleware)

	projectsGroup.POST("", BucketController.Store)
	projectsGroup.GET("", BucketController.List)
	projectsGroup.GET("/:bucketUUID", BucketController.Show)
	projectsGroup.PUT("/:bucketUUID", BucketController.Update)
	projectsGroup.DELETE("/:bucketUUID", BucketController.Delete)

	filesGroup := projectsGroup.Group("/:bucketUUID/files")

	filesGroup.POST("", FileController.Store)
	filesGroup.GET("", FileController.List)
	filesGroup.GET("/:fileUUID", FileController.Show)
	filesGroup.PUT("/:fileUUID", FileController.Rename)
	filesGroup.DELETE("/:fileUUID", FileController.Delete)
}
