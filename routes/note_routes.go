package routes

import (
	"github.com/labstack/echo/v4"
	"myapp/controllers"
	"myapp/middleware"
)

func RegisterNoteRoutes(e *echo.Echo, noteController *controllers.NoteController) {
	notesGroup := e.Group("/notes", middleware.AuthMiddleware)

	notesGroup.POST("", noteController.Store)
	notesGroup.GET("", noteController.List)
	notesGroup.GET("/:id", noteController.Show)
	notesGroup.PUT("/:id", noteController.Update)
	notesGroup.DELETE("/:id", noteController.Delete)
}
