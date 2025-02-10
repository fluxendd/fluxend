package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc, userController *controllers.UserController) {
	e.POST("api/users/register", userController.Store)
	e.POST("api/users/login", userController.Login)
	// e.GET("/users", userController.List)
	e.GET("/api/users/:id", authMiddleware(userController.Show))
	e.PUT("/api/users/:id", authMiddleware(userController.Update))
	// e.DELETE("/users/:id", userController.Delete)
}
