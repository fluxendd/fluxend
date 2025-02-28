package controllers

import (
	"fluxton/requests/form_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormFieldController struct {
	formFieldService services.FormFieldService
}

func NewFormFieldController(injector *do.Injector) (*FormFieldController, error) {
	formFieldService := do.MustInvoke[services.FormFieldService](injector)

	return &FormFieldController{formFieldService: formFieldService}, nil
}

func (ffc *FormFieldController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.List(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResourceCollection(formFields))
}

func (ffc *FormFieldController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formField, err := ffc.formFieldService.GetByUUID(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResource(&formField))
}

func (ffc *FormFieldController) Store(c echo.Context) error {
	var request form_requests.CreateFieldRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formField, err := ffc.formFieldService.Create(formUUID, projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormFieldResource(&formField))
}

func (ffc *FormFieldController) Update(c echo.Context) error {
	var request form_requests.CreateFieldRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formFieldUUID, err := utils.GetUUIDPathParam(c, "formFieldUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedFormField, err := ffc.formFieldService.Update(formUUID, formFieldUUID, projectUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResource(updatedFormField))
}

func (ffc *FormFieldController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	formFieldUUID, err := utils.GetUUIDPathParam(c, "formFieldUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := ffc.formFieldService.Delete(formUUID, formFieldUUID, projectUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (ffc *FormFieldController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, error) {
	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return projectUUID, formUUID, nil
}
