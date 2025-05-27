package routes

import (
	"fluxend/internal/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

func RegisterFormRoutes(
	e *echo.Echo,
	container *do.Injector,
	authMiddleware echo.MiddlewareFunc,
	formEnabledMiddleware echo.MiddlewareFunc,
) {
	formController := do.MustInvoke[*handlers.FormHandler](container)
	formFieldController := do.MustInvoke[*handlers.FormFieldHandler](container)
	formResponseController := do.MustInvoke[*handlers.FormResponseHandler](container)

	formsGroup := e.Group("forms", authMiddleware, formEnabledMiddleware)

	formsGroup.POST("", formController.Store)
	formsGroup.GET("", formController.List)
	formsGroup.GET("/:formUUID", formController.Show)
	formsGroup.PUT("/:formUUID", formController.Update)
	formsGroup.DELETE("/:formUUID", formController.Delete)

	// Form Field routes
	formFieldsGroup := e.Group("api/forms/:formUUID/fields", authMiddleware, formEnabledMiddleware)

	formFieldsGroup.POST("", formFieldController.Store)
	formFieldsGroup.GET("", formFieldController.List)
	formFieldsGroup.GET("/:fieldUUID", formFieldController.Show)
	formFieldsGroup.PUT("/:fieldUUID", formFieldController.Update)
	formFieldsGroup.DELETE("/:fieldUUID", formFieldController.Delete)

	// Form Response routes
	formResponsesGroup := e.Group("api/forms/:formUUID/responses", authMiddleware, formEnabledMiddleware)

	formResponsesGroup.GET("", formResponseController.List)
	formResponsesGroup.POST("", formResponseController.Store)
	formResponsesGroup.GET("/:formResponseUUID", formResponseController.Show)
	formResponsesGroup.DELETE("/:formResponseUUID", formResponseController.Delete)
}
