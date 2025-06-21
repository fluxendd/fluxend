package handlers

import (
	"fluxend/internal/api/dto"
	organizationDto "fluxend/internal/api/dto/organization"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/organization"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type OrganizationHandler struct {
	organizationService organization.Service
}

func NewOrganizationHandler(injector *do.Injector) (*OrganizationHandler, error) {
	organizationService := do.MustInvoke[organization.Service](injector)

	return &OrganizationHandler{organizationService: organizationService}, nil
}

// List all organizations
//
// @Summary List organizations
// @Description Get all organizations
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} response.Response{content=[]organization.Response} "List of organizations"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations [get]
func (oh *OrganizationHandler) List(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUserId, _ := auth.NewAuth(c).Uuid()

	paginationParams := request.ExtractPaginationParams(c)
	organizations, err := oh.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToOrganizationResourceCollection(organizations))
}

// Show details of a single organization
//
// @Summary Retrieve organization
// @Description Get details of a specific organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
//
// @Success 200 {object} response.Response{content=organization.Response} "Organization details"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /organizations/{organizationUUID} [get]
func (oh *OrganizationHandler) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedOrganization, err := oh.organizationService.GetByID(organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToOrganizationResource(&fetchedOrganization))
}

// Store creates a new organization
//
// @Summary Create organization
// @Description Add a new organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization body organization.CreateRequest true "Organization name"
//
// @Success 201 {object} response.Response{content=organization.Response} "Organization created"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /organizations [post]
func (oh *OrganizationHandler) Store(c echo.Context) error {
	var request organizationDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, err := auth.NewAuth(c).User()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	storedOrganization, err := oh.organizationService.Create(request.Name, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToOrganizationResource(&storedOrganization))
}

// Update an organization
//
// @Summary Update organization
// @Description Modify an existing organization's details
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
// @Param organization body organization.CreateRequest true "Updated organization details"
//
// @Success 200 {object} response.Response{content=organization.Response} "Organization updated"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /organizations/{organizationUUID} [put]
func (oh *OrganizationHandler) Update(c echo.Context) error {
	var request organizationDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := oh.organizationService.Update(request.Name, organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToOrganizationResource(updatedOrganization))
}

// Delete an organization
//
// @Summary Delete organization
// @Description Remove an organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
//
// @Success 204 {object} response.Response{} "Organization deleted"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [delete]
func (oh *OrganizationHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := oh.organizationService.Delete(organizationUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
