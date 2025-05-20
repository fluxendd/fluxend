package handlers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/form"
	formMapper "fluxton/internal/api/mapper/form"
	"fluxton/internal/api/response"
	formDomain "fluxton/internal/domain/form"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FormFieldHandler struct {
	formFieldService formDomain.FieldService
}

func NewFormFieldHandler(injector *do.Injector) (*FormFieldHandler, error) {
	formFieldService := do.MustInvoke[formDomain.FieldService](injector)

	return &FormFieldHandler{formFieldService: formFieldService}, nil
}

// List retrieves all fields for a specific form
//
// @Summary List all fields for a form
// @Description Retrieve a list of all fields in a specific form
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
//
// @Success 200 {array} responses.Response{content=[]resources.FieldResponse} "List of fields"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields [get]
func (ffc *FormFieldHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.List(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formMapper.ToFieldResourceCollection(formFields))
}

// Show retrieves details of a specific field
//
// @Summary Show details of a single field
// @Description Get details of a specific field
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 200 {object} responses.Response{content=resources.FieldResponse} "Field details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [get]
func (ffc *FormFieldHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formField, err := ffc.formFieldService.GetByUUID(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formMapper.ToFieldResource(&formField))
}

// Store creates a new field for a form
//
// @Summary Create a new field for a form
// @Description Add a new field to a form
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.CreateFormFieldsRequest true "Field details"
// @Param formUUID path string true "Form UUID"
//
// @Success 201 {object} responses.Response{content=resources.FieldResponse} "Field created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields [post]
func (ffc *FormFieldHandler) Store(c echo.Context) error {
	var request form.CreateFormFieldsRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffc.formFieldService.CreateMany(formUUID, form.ToCreateFormFieldInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, formMapper.ToFieldResourceCollection(formFields))
}

// Update updates an existing field
//
// @Summary Update an existing field
// @Description Update the details of an existing field
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form_requests.UpdateFormFieldRequest true "Field details"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 200 {object} responses.Response{content=resources.FieldResponse} "Field updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [put]
func (ffc *FormFieldHandler) Update(c echo.Context) error {
	var request form.UpdateFormFieldRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fieldUUID, err := request.GetUUIDPathParam(c, "fieldUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedFormField, err := ffc.formFieldService.Update(formUUID, fieldUUID, authUser, form.ToUpdateFormFieldInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, formMapper.ToFieldResource(updatedFormField))
}

// Delete deletes a field from a form
//
// @Summary Delete a field from a form
// @Description Remove a specific field from the form
// @Tags Form Fields
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
func (ffc *FormFieldHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fieldUUID, err := request.GetUUIDPathParam(c, "fieldUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := ffc.formFieldService.Delete(formUUID, fieldUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
