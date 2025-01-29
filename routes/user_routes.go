package routes

import (
	"github.com/labstack/echo/v4"
	"myapp/controllers"
	"myapp/middleware"
)

func RegisterUserRoutes(e *echo.Echo, userController *controllers.UserController) {
	e.POST("/users/register", userController.Store)
	e.POST("/users/login", userController.Login)
	// e.GET("/users", userController.List)
	e.GET("/users/:id", middleware.AuthMiddleware(userController.Show))
	e.PUT("/users/:id", middleware.AuthMiddleware(userController.Update))
	// e.DELETE("/users/:id", userController.Delete)
}
