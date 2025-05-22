package handlers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/form"
	formMapper "fluxton/internal/api/mapper/form"
	"fluxton/internal/api/response"
	formDomain "fluxton/internal/domain/form"
	"fluxton/pkg/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormResponseHandler struct {
	formResponseService formDomain.FieldResponseService
}

func NewFormResponseHandler(injector *do.Injector) (*FormResponseHandler, error) {
	formResponseService := do.MustInvoke[formDomain.FieldResponseService](injector)

	return &FormResponseHandler{formResponseService: formResponseService}, nil
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
// @Success 200 {object} response.Response{content=[]form.ResponseForAPI} "List of form responses"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses [get]
func (ffh *FormResponseHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formResponses, err := ffh.formResponseService.List(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formMapper.ToResponseResourceCollection(formResponses))
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
// @Success 200 {object} response.Response{content=form.ResponseForAPI} "Form response details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses/{formResponseUUID} [get]
func (ffh *FormResponseHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, formResponseUUID, err := ffh.parseRequest(request, c)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffh.formResponseService.GetByUUID(formResponseUUID, formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formMapper.ToResponseResource(formResponse))
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
// @Param request body form.CreateResponseRequest true "Request body to create a new form response"
//
// @Success 201 {object} response.Response{content=form.ResponseForAPI} "Form response details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/responses [post]
func (ffh *FormResponseHandler) Store(c echo.Context) error {
	var request form.CreateResponseRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formResponse, err := ffh.formResponseService.Create(formUUID, form.ToCreateFormResponseInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, formMapper.ToResponseResource(&formResponse))
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
func (ffh *FormResponseHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, formResponseUUID, err := ffh.parseRequest(request, c)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if err = ffh.formResponseService.Delete(formUUID, formResponseUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}

func (ffh *FormResponseHandler) parseRequest(request dto.DefaultRequestWithProjectHeader, c echo.Context) (uuid.UUID, uuid.UUID, error) {
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
