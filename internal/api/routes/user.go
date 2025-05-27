package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterUserRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	userController := do.MustInvoke[*handlers.UserHandler](container)

	e.POST("users/register", userController.Store)
	e.POST("users/login", userController.Login)
	e.GET("users/:userUUID", authMiddleware(userController.Show))
	e.PUT("users/:userUUID", authMiddleware(userController.Update))
	e.POST("users/logout", authMiddleware(userController.Logout))
}
