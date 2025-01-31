package routes

import (
	"fluxton/controllers"
	"fluxton/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterNoteRoutes(e *echo.Echo, noteController *controllers.NoteController) {
	notesGroup := e.Group("api/notes", middleware.AuthMiddleware)

	notesGroup.POST("", noteController.Store)
	notesGroup.GET("", noteController.List)
	notesGroup.GET("/:id", noteController.Show)
	notesGroup.PUT("/:id", noteController.Update)
	notesGroup.DELETE("/:id", noteController.Delete)
}
