package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterFunctionRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	functionController := do.MustInvoke[*handlers.FunctionHandler](container)

	functionsGroup := e.Group("functions", authMiddleware)

	// table routes
	functionsGroup.GET("/:schema", functionController.List)
	functionsGroup.POST("/:schema", functionController.Store)
	functionsGroup.GET("/:schema/:functionName", functionController.Show)
	functionsGroup.DELETE("/:schema/:functionName", functionController.Delete)
}
