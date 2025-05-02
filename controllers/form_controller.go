package controllers

import (
	"fluxton/requests"
	"fluxton/requests/form_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormController struct {
	formService services.FormService
}

func NewFormController(injector *do.Injector) (*FormController, error) {
	formService := do.MustInvoke[services.FormService](injector)

	return &FormController{formService: formService}, nil
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
func (fc *FormController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}
	authUser, _ := utils.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	forms, err := fc.formService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResourceCollection(forms))
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
func (fc *FormController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	form, err := fc.formService.GetByUUID(formUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(&form))
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
func (fc *FormController) Store(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	form, err := fc.formService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormResource(&form))
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
func (fc *FormController) Update(c echo.Context) error {
	var request form_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedForm, err := fc.formService.Update(formUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormResource(updatedForm))
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
func (fc *FormController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := fc.formService.Delete(formUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
