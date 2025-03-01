package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterBucketRoutes(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	BucketController *controllers.BucketController,
) {
	projectsGroup := e.Group("api/projects/:projectUUID/buckets", authMiddleware)

	projectsGroup.POST("", BucketController.Store)
	projectsGroup.GET("", BucketController.List)
	projectsGroup.GET("/:bucketUUID", BucketController.Show)
	projectsGroup.PUT("/:bucketUUID", BucketController.Update)
	projectsGroup.DELETE("/:bucketUUID", BucketController.Delete)
}
