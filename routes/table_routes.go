package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterTableRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	tableController := do.MustInvoke[*controllers.TableController](container)
	columnController := do.MustInvoke[*controllers.ColumnController](container)
	indexController := do.MustInvoke[*controllers.IndexController](container)

	tablesGroup := e.Group("api/projects/:projectID/tables", authMiddleware)

	// table routes
	tablesGroup.POST("", tableController.Store)
	tablesGroup.GET("", tableController.List)
	tablesGroup.GET("/:tableID", tableController.Show)
	tablesGroup.PUT("/:tableID/duplicate", tableController.Duplicate)
	tablesGroup.PUT("/:tableID/rename", tableController.Rename)
	tablesGroup.DELETE("/:tableID", tableController.Delete)

	// column routes
	tablesGroup.POST("/:tableID/columns", columnController.Store)
	tablesGroup.PUT("/:tableID/columns", columnController.Alter)
	tablesGroup.PUT("/:tableID/columns/:columnName", columnController.Rename)
	tablesGroup.DELETE("/:tableID/columns/:columnName", columnController.Delete)

	// index routes
	tablesGroup.POST("/:tableID/indexes", indexController.Store)
	tablesGroup.GET("/:tableID/indexes", indexController.List)
	tablesGroup.GET("/:tableID/indexes/:indexName", indexController.Show)
	tablesGroup.DELETE("/:tableID/indexes/:indexName", indexController.Delete)
}
