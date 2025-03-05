package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterFormRoutes(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	FormController *controllers.FormController,
	FormFieldController *controllers.FormFieldController,
	FormResponseController *controllers.FormResponseController,
) {
	projectsGroup := e.Group("api/projects/:projectUUID/forms", authMiddleware)

	projectsGroup.POST("", FormController.Store)
	projectsGroup.GET("", FormController.List)
	projectsGroup.GET("/:formUUID", FormController.Show)
	projectsGroup.PUT("/:formUUID", FormController.Update)
	projectsGroup.DELETE("/:formUUID", FormController.Delete)

	// Form Field routes
	formFieldsGroup := e.Group("api/projects/:projectUUID/forms/:formUUID/fields", authMiddleware)

	formFieldsGroup.POST("", FormFieldController.Store)
	formFieldsGroup.GET("", FormFieldController.List)
	formFieldsGroup.GET("/:fieldUUID", FormFieldController.Show)
	formFieldsGroup.PUT("/:fieldUUID", FormFieldController.Update)
	formFieldsGroup.DELETE("/:fieldUUID", FormFieldController.Delete)

	// Form Response routes
	formResponsesGroup := e.Group("api/projects/:projectUUID/forms/:formUUID/responses", authMiddleware)

	formResponsesGroup.GET("", FormResponseController.List)
	formResponsesGroup.POST("", FormResponseController.Store)
}
