package routes

import (
	"fluxton/controllers"
	"fluxton/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRowRoutes(e *echo.Echo, RowController *controllers.RowController) {
	rowsGroup := e.Group("api/c/:projectID/:tableName", middleware.AuthMiddleware)

	// table routes
	rowsGroup.POST("", RowController.Store)
	rowsGroup.GET("", RowController.List)
	rowsGroup.GET("/:rowID", RowController.Show)
}
