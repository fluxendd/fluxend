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

type FormResponseController struct {
	formResponseService services.FormResponseService
}

func NewFormResponseController(injector *do.Injector) (*FormResponseController, error) {
	formResponseService := do.MustInvoke[services.FormResponseService](injector)

	return &FormResponseController{formResponseService: formResponseService}, nil
}

func (ffc *FormResponseController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponses, err := ffc.formResponseService.List(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResponseResourceCollection(formResponses))
}

func (ffc *FormResponseController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, formResponseUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffc.formResponseService.GetByUUID(formResponseUUID, formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResponseResource(formResponse))
}

func (ffc *FormResponseController) Store(c echo.Context) error {
	var request form_requests.CreateResponseRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffc.formResponseService.Create(formUUID, projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResponseResource(&formResponse))
}

func (ffc *FormResponseController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, formResponseUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err = ffc.formResponseService.Delete(projectUUID, formUUID, formResponseUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (ffc *FormResponseController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	formResponseUUID, err := utils.GetUUIDPathParam(c, "formResponseUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, err
	}

	return projectUUID, formUUID, formResponseUUID, nil
}
