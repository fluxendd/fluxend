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

type FormFieldHandler struct {
	formFieldService form.FieldService
}

func NewFormFieldHandler(injector *do.Injector) (*FormFieldHandler, error) {
	formFieldService := do.MustInvoke[form.FieldService](injector)

	return &FormFieldHandler{formFieldService: formFieldService}, nil
}

// List retrieves all fields for a specific form
//
// @Summary List fields
// @Description Retrieve a list of all fields in a specific form
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param formUUID path string true "Form UUID"
//
// @Success 200 {array} response.Response{content=[]form.FieldResponseApi} "List of fields"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID}/fields [get]
func (ffh *FormFieldHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffh.formFieldService.List(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFieldResourceCollection(formFields))
}

// Show retrieves details of a specific field
//
// @Summary Retrieve field
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
// @Success 200 {object} response.Response{content=form.FieldResponseApi} "Field details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [get]
func (ffh *FormFieldHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formField, err := ffh.formFieldService.GetByUUID(formUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFieldResource(&formField))
}

// Store creates a new field for a form
//
// @Summary Create field
// @Description Add a new field to a form
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form.CreateFormFieldsRequest true "Field details"
// @Param formUUID path string true "Form UUID"
//
// @Success 201 {object} response.Response{content=form.FieldResponseApi} "Field created"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID}/fields [post]
func (ffh *FormFieldHandler) Store(c echo.Context) error {
	var request formDto.CreateFormFieldsRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	formUUID, err := request.GetUUIDPathParam(c, "formUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	formFields, err := ffh.formFieldService.CreateMany(formUUID, formDto.ToCreateFormFieldInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToFieldResourceCollection(formFields))
}

// Update updates an existing field
//
// @Summary Update field
// @Description Update the details of an existing field
// @Tags Form Fields
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param field body form.UpdateFormFieldRequest true "Field details"
// @Param formUUID path string true "Form UUID"
// @Param fieldUUID path string true "Field UUID"
//
// @Success 200 {object} response.Response{content=form.FieldResponseApi} "Field updated"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [put]
func (ffh *FormFieldHandler) Update(c echo.Context) error {
	var request formDto.UpdateFormFieldRequest
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

	updatedFormField, err := ffh.formFieldService.Update(formUUID, fieldUUID, authUser, formDto.ToUpdateFormFieldInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFieldResource(updatedFormField))
}

// Delete deletes a field from a form
//
// @Summary Delete field
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
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /forms/{formUUID}/fields/{fieldUUID} [delete]
func (ffh *FormFieldHandler) Delete(c echo.Context) error {
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

	if _, err := ffh.formFieldService.Delete(formUUID, fieldUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
