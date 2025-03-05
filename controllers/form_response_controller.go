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
		return err
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return err
	}

	formResponses, err := ffc.formResponseService.List(formUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResponseResourceCollection(formResponses))
}

func (ffc *FormResponseController) Store(c echo.Context) error {
	var request form_requests.CreateResponseRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return err
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return err
	}

	formResponse, err := ffc.formResponseService.Create(formUUID, projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResponseResource(&formResponse))
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
