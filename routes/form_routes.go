package routes

import (
	"fluxton/controllers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterFormRoutes(e *echo.Echo, container *do.Injector, authMiddleware echo.MiddlewareFunc) {
	formController := do.MustInvoke[*controllers.FormController](container)
	formFieldController := do.MustInvoke[*controllers.FormFieldController](container)
	formResponseController := do.MustInvoke[*controllers.FormResponseController](container)

	formsGroup := e.Group("api/forms", authMiddleware)

	formsGroup.POST("", formController.Store)
	formsGroup.GET("", formController.List)
	formsGroup.GET("/:formUUID", formController.Show)
	formsGroup.PUT("/:formUUID", formController.Update)
	formsGroup.DELETE("/:formUUID", formController.Delete)

	// Form Field routes
	formFieldsGroup := e.Group("api/forms/:formUUID/fields", authMiddleware)

	formFieldsGroup.POST("", formFieldController.Store)
	formFieldsGroup.GET("", formFieldController.List)
	formFieldsGroup.GET("/:fieldUUID", formFieldController.Show)
	formFieldsGroup.PUT("/:fieldUUID", formFieldController.Update)
	formFieldsGroup.DELETE("/:fieldUUID", formFieldController.Delete)

	// Form Response routes
	formResponsesGroup := e.Group("api/forms/:formUUID/responses", authMiddleware)

	formResponsesGroup.GET("", formResponseController.List)
	formResponsesGroup.POST("", formResponseController.Store)
	formResponsesGroup.GET("/:formResponseUUID", formResponseController.Show)
	formResponsesGroup.DELETE("/:formResponseUUID", formResponseController.Delete)
}
