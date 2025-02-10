package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterRowRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, RowController *controllers.RowController) {
	rowsGroup := e.Group("api/c/:projectID/:tableName", authMiddleware)

	// table routes
	rowsGroup.POST("", RowController.Store)
	rowsGroup.GET("", RowController.List)
	rowsGroup.GET("/:rowID", RowController.Show)
}
