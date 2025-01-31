package routes

import (
	"fluxton/controllers"
	"fluxton/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterTableRoutes(e *echo.Echo, TableController *controllers.TableController) {
	tablesGroup := e.Group("api/projects/:projectID/tables", middleware.AuthMiddleware)

	tablesGroup.POST("", TableController.Store)
	tablesGroup.GET("", TableController.List)
	tablesGroup.GET("/:tableID", TableController.Show)
	/*
		tablesGroup.PUT("/:tableID", TableController.Update)
		tablesGroup.DELETE("/:tableID", TableController.Delete)*/
}
