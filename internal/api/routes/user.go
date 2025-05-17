package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterUserRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	userController := do.MustInvoke[*controllers.UserController](container)

	e.POST("api/users/register", userController.Store)
	e.POST("api/users/login", userController.Login)
	e.GET("/api/users/:userUUID", authMiddleware(userController.Show))
	e.PUT("/api/users/:userUUID", authMiddleware(userController.Update))
	e.POST("api/users/logout", authMiddleware(userController.Logout))
}
