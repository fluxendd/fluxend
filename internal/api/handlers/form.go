package handlers

import (
	"fluxton/internal/api/dto"
	formDto "fluxton/internal/api/dto/form"
	"fluxton/internal/api/response"
	formDomain "fluxton/internal/domain/form"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormHandler struct {
	formService formDomain.Service
}

func NewFormHandler(injector *do.Injector) (*FormHandler, error) {
	formService := do.MustInvoke[formDomain.Service](injector)

	return &FormHandler{formService: formService}, nil
}

// List retrieves all forms for a project
//
// @Summary List all forms
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
// @Success 200 {array} responses.Response{content=[]resources.FormResponse} "List of forms"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms [get]
func (fc *FormHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}
	authUser, _ := auth.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	forms, err := fc.formService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formDto.FormResourceCollection(forms))
}

// Show retrieves details of a specific form
//
// @Summary Show details of a single form
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
// @Success 200 {object} responses.Response{content=resources.FormResponse} "Form details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID} [get]
func (fc *FormHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	form, err := fc.formService.GetByUUID(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, form.FormResource(&form))
}

// Store creates a new form
//
// @Summary Create a new form
// @Description Add a new form with a name and description
// @Tags Forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param form body form_requests.CreateRequest true "Form name and description"
//
// @Success 201 {object} responses.Response{content=resources.FormResponse} "Form created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms [post]
func (fc *FormHandler) Store(c echo.Context) error {
	var request form.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	form, err := fc.formService.Create(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, form.FormResource(&form))
}

// Update updates an existing form
//
// @Summary Update an existing form
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
// @Param form body form_requests.CreateRequest true "Form name and description"
//
// @Success 200 {object} responses.Response{content=resources.FormResponse} "Form updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID} [put]
func (fc *FormHandler) Update(c echo.Context) error {
	var request form.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedForm, err := fc.formService.Update(formUUID, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, form.FormResource(updatedForm))
}

// Delete removes a form
//
// @Summary Delete a form
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
func (fc *FormHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := fc.formService.Delete(formUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
