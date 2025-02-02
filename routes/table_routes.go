package routes

import (
	"fluxton/controllers"
	"fluxton/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterTableRoutes(e *echo.Echo, TableController *controllers.TableController, ColumnController *controllers.ColumnController) {
	tablesGroup := e.Group("api/projects/:projectID/tables", middleware.AuthMiddleware)

	// table routes
	tablesGroup.POST("", TableController.Store)
	tablesGroup.GET("", TableController.List)
	tablesGroup.GET("/:tableID", TableController.Show)
	tablesGroup.PUT("/:tableID", TableController.Rename)
	tablesGroup.DELETE("/:tableID", TableController.Delete)

	// column routes
	tablesGroup.POST("/:tableID/columns", ColumnController.Store)
	tablesGroup.PUT("/:tableID/columns/:columnName", ColumnController.Alter)
	tablesGroup.DELETE("/:tableID/columns/:columnName", ColumnController.Delete)
}
