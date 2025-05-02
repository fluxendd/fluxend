package controllers

import (
	"fluxton/requests"
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

// List all form responses for a form
//
// @Summary List all form responses for a form
// @Description Get all form responses for a specific form
// @Tags Form Responses
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param formUUID path string true "Form UUID"
//
// @Success 200 {object} responses.Response{content=[]resources.FormResponseForAPI} "List of form responses"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses [get]
func (ffc *FormResponseController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponses, err := ffc.formResponseService.List(formUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResponseResourceCollection(formResponses))
}

// Show details of a single form response
//
// @Summary Show details of a single form response
// @Description Get details of a specific form response
// @Tags Form Responses
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param projectUUID path string true "Project UUID"
// @Param formUUID path string true "Form UUID"
// @Param formResponseUUID path string true "Form Response UUID"
//
// @Success 200 {object} responses.Response{content=resources.FormResponseForAPI} "Form response details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses/{formResponseUUID} [get]
func (ffc *FormResponseController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, formResponseUUID, err := ffc.parseRequest(request, c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffc.formResponseService.GetByUUID(formResponseUUID, formUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResponseResource(formResponse))
}

// Store a new form response
//
// @Summary Store a new form response
// @Description Store a new form response for a specific form
// @Tags Form Responses
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param formUUID path string true "Form UUID"
//
// @Param request body form_requests.CreateResponseRequest true "Request body to create a new form response"
//
// @Success 201 {object} responses.Response{content=resources.FormResponseForAPI} "Form response details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses [post]
func (ffc *FormResponseController) Store(c echo.Context) error {
	var request form_requests.CreateResponseRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffc.formResponseService.Create(formUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResponseResource(&formResponse))
}

// Delete a form response
//
// @Summary Delete a form response
// @Description Delete a specific form response
// @Tags Form Responses
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param projectUUID path string true "Project UUID"
// @Param formUUID path string true "Form UUID"
// @Param formResponseUUID path string true "Form Response UUID"
//
// @Success 204 "Form response deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses/{formResponseUUID} [delete]
func (ffc *FormResponseController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, formResponseUUID, err := ffc.parseRequest(request, c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err = ffc.formResponseService.Delete(formUUID, formResponseUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (ffc *FormResponseController) parseRequest(request requests.DefaultRequestWithProjectHeader, c echo.Context) (uuid.UUID, uuid.UUID, error) {
	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	formResponseUUID, err := request.GetUUIDPathParam(c, "formResponseUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return formUUID, formResponseUUID, nil
}
