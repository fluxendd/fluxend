package routes

import (
	"fluxton/controllers"
	"fluxton/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, userController *controllers.UserController) {
	e.POST("api/users/register", userController.Store)
	e.POST("api/users/login", userController.Login)
	// e.GET("/users", userController.List)
	e.GET("/api/users/:id", middlewares.AuthMiddleware(userController.Show))
	e.PUT("/api/users/:id", middlewares.AuthMiddleware(userController.Update))
	// e.DELETE("/users/:id", userController.Delete)
}
