package controllers

import (
	"fluxton/requests/form_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
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
// @Tags Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
//
// @Success 200 {array} responses.Response{content=[]resources.FormFieldResponse} "List of fields"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields [get]
func (ffc *FormFieldController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.List(formUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResourceCollection(formFields))
}

// Show retrieves details of a specific field
//
// @Summary Show details of a single field
// @Description Get details of a specific field
// @Tags Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 200 {object} responses.Response{content=resources.FormFieldResponse} "Field details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [get]
func (ffc *FormFieldController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formField, err := ffc.formFieldService.GetByUUID(formUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResource(&formField))
}

// Store creates a new field for a form
//
// @Summary Create a new field for a form
// @Description Add a new field to a form
// @Tags Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.CreateFormFieldsRequest true "Field details"
// @Param formUUID path string true "Form UUID"
//
// @Success 201 {object} responses.Response{content=resources.FormFieldResponse} "Field created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields [post]
func (ffc *FormFieldController) Store(c echo.Context) error {
	var request form_requests.CreateFormFieldsRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.CreateMany(formUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FormFieldResourceCollection(formFields))
}

// Update updates an existing field
//
// @Summary Update an existing field
// @Description Update the details of an existing field
// @Tags Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.UpdateFormFieldRequest true "Field details"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 200 {object} responses.Response{content=resources.FormFieldResponse} "Field updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [put]
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

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedFormField, err := ffc.formFieldService.Update(formUUID, fieldUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FormFieldResource(updatedFormField))
}

// Delete deletes a field from a form
//
// @Summary Delete a field from a form
// @Description Remove a specific field from the form
// @Tags Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 204 "Field deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [delete]
func (ffc *FormFieldController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	fieldUUID, err := utils.GetUUIDPathParam(c, "fieldUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	formUUID, err := utils.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := ffc.formFieldService.Delete(formUUID, fieldUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
