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

	tablesGroup := e.Group("api/tables", authMiddleware)

	// table routes
	tablesGroup.POST("", tableController.Store)
	tablesGroup.GET("", tableController.List)
	tablesGroup.GET("/:tableUUID", tableController.Show)
	tablesGroup.PUT("/:tableUUID/duplicate", tableController.Duplicate)
	tablesGroup.PUT("/:tableUUID/rename", tableController.Rename)
	tablesGroup.DELETE("/:tableUUID", tableController.Delete)

	// column routes
	tablesGroup.POST("/:tableUUID/columns", columnController.Store)
	tablesGroup.PUT("/:tableUUID/columns", columnController.Alter)
	tablesGroup.PUT("/:tableUUID/columns/:columnName", columnController.Rename)
	tablesGroup.DELETE("/:tableUUID/columns/:columnName", columnController.Delete)

	// index routes
	tablesGroup.POST("/:tableUUID/indexes", indexController.Store)
	tablesGroup.GET("/:tableUUID/indexes", indexController.List)
	tablesGroup.GET("/:tableUUID/indexes/:indexName", indexController.Show)
	tablesGroup.DELETE("/:tableUUID/indexes/:indexName", indexController.Delete)
}
