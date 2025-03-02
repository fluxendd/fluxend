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

// List retrieves all fields for a specific form
//
// @Summary List all fields for a form
// @Description Retrieve a list of all fields in a specific form
// @Tags fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formId path string true "Form ID"
//
// @Success 200 {array} responses.Response{content=[]resources.FormFieldResponse} "List of fields"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formId}/fields [get]
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

// Store creates a new field for a form
//
// @Summary Create a new field for a form
// @Description Add a new field to a form
// @Tags fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.CreateFormFieldsRequest true "Field details"
// @Param formId path string true "Form ID"
//
// @Success 201 {object} responses.Response{content=resources.FormFieldResponse} "Field created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formId}/fields [post]
func (ffc *FormFieldController) Store(c echo.Context) error {
	var request form_requests.CreateFormFieldsRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.CreateMany(formUUID, projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormFieldResourceCollection(formFields))
}

// Update updates an existing field
//
// @Summary Update an existing field
// @Description Update the details of an existing field
// @Tags fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.UpdateFormFieldRequest true "Field details"
// @Param formId path string true "Form ID"
// @Param fieldId path string true "Field ID"
//
// @Success 200 {object} responses.Response{content=resources.FormFieldResponse} "Field updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formId}/fields/{fieldId} [put]
func (ffc *FormFieldController) Update(c echo.Context) error {
	var request form_requests.UpdateFormFieldRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fieldUUID, err := utils.GetUUIDPathParam(c, "fieldUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedFormField, err := ffc.formFieldService.Update(formUUID, fieldUUID, projectUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResource(updatedFormField))
}

// Delete deletes a field from a form
//
// @Summary Delete a field from a form
// @Description Remove a specific field from the form
// @Tags fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formId path string true "Form ID"
// @Param fieldId path string true "Field ID"
//
// @Success 204 "Field deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formId}/fields/{fieldId} [delete]
func (ffc *FormFieldController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	fieldUUID, err := utils.GetUUIDPathParam(c, "fieldUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	projectUUID, formUUID, err := ffc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := ffc.formFieldService.Delete(formUUID, fieldUUID, projectUUID, authUser); err != nil {
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
