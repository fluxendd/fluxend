package handlers

import (
	"fluxend/internal/api/dto"
	formDto "fluxend/internal/api/dto/form"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/form"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormHandler struct {
	formService form.Service
}

func NewFormHandler(injector *do.Injector) (*FormHandler, error) {
	formService := do.MustInvoke[form.Service](injector)

	return &FormHandler{formService: formService}, nil
}

// List retrieves all forms for a project
//
// @Summary List forms
// @Description Retrieve a list of all forms for the specified project
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {array} response.Response{content=[]form.Response} "List of forms"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms [get]
func (fh *FormHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}
	authUser, _ := auth.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	forms, err := fh.formService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFormResourceCollection(forms))
}

// Show retrieves details of a specific form
//
// @Summary Retrieve form
// @Description Get details of a specific form
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param formUUID path string true "Form UUID"
//
// @Success 200 {object} response.Response{content=form.Response} "Form details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID} [get]
func (fh *FormHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedForm, err := fh.formService.GetByUUID(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFormResource(&fetchedForm))
}

// Store creates a new form
//
// @Summary Create form
// @Description Add a new form with a name and description
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param form body form.CreateRequest true "Form name and description"
//
// @Success 201 {object} response.Response{content=form.Response} "Form created"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms [post]
func (fh *FormHandler) Store(c echo.Context) error {
	var request formDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	createdForm, err := fh.formService.Create(formDto.ToCreateFormInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToFormResource(&createdForm))
}

// Update updates an existing form
//
// @Summary Update form
// @Description Update form details such as name and description
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param formUUID path string true "Form UUID"
// @Param form body form.CreateRequest true "Form name and description"
//
// @Success 200 {object} response.Response{content=form.Response} "Form updated"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID} [put]
func (fh *FormHandler) Update(c echo.Context) error {
	var request formDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedForm, err := fh.formService.Update(formUUID, authUser, formDto.ToCreateFormInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFormResource(updatedForm))
}

// Delete removes a form
//
// @Summary Delete form
// @Description Remove a form from the project
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param formUUID path string true "Form UUID"
//
// @Success 204 "Form deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID} [delete]
func (fh *FormHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := fh.formService.Delete(formUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
