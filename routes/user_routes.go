package routes

import (
	"github.com/labstack/echo/v4"
	"myapp/controllers"
	"myapp/middleware"
)

func RegisterUserRoutes(e *echo.Echo, userController *controllers.UserController) {
	e.POST("api/users/register", userController.Store)
	e.POST("api/users/login", userController.Login)
	// e.GET("/users", userController.List)
	e.GET("/api/users/:id", middleware.AuthMiddleware(userController.Show))
	e.PUT("/api/users/:id", middleware.AuthMiddleware(userController.Update))
	// e.DELETE("/users/:id", userController.Delete)
}
