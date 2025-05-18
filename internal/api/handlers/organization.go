package handlers

import (
	"fluxton/internal/api/dto"
	organizationDto "fluxton/internal/api/dto/organization"
	organizationMapper "fluxton/internal/api/mapper/organization"
	"fluxton/internal/api/response"
	organizationDomain "fluxton/internal/domain/organization"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type OrganizationHandler struct {
	organizationService organizationDomain.Service
}

func NewOrganizationHandler(injector *do.Injector) (*OrganizationHandler, error) {
	organizationService := do.MustInvoke[organizationDomain.Service](injector)

	return &OrganizationHandler{organizationService: organizationService}, nil
}

// List all organizations
//
// @Summary List all organizations
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
// @Success 200 {object} responses.Response{content=[]organizationDto.Response} "List of organizations"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations [get]
func (oc *OrganizationHandler) List(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUserId, _ := auth.NewAuth(c).Uuid()

	paginationParams := request.ExtractPaginationParams(c)
	organizations, err := oc.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, organizationMapper.ToResourceCollection(organizations))
}

// Show details of a single organization
//
// @Summary Show details of a single organization
// @Description Get details of a specific organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
//
// @Success 200 {object} responses.Response{content=organizationDto.Response} "Organization details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [get]
func (oc *OrganizationHandler) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	organization, err := oc.organizationService.GetByID(organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, organizationMapper.ToResource(&organization))
}

// Store creates a new organization
//
// @Summary Create a new organization
// @Description Add a new organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization body organization_requests.CreateRequest true "Organization name"
//
// @Success 201 {object} responses.Response{content=organizationDto.Response} "Organization created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations [post]
func (oc *OrganizationHandler) Store(c echo.Context) error {
	var request organizationDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, err := auth.NewAuth(c).User()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	organization, err := oc.organizationService.Create(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, organizationMapper.ToResource(&organization))
}

// Update an organization
//
// @Summary Update an organization
// @Description Modify an existing organization's details
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
// @Param organization body organization_requests.CreateRequest true "Updated organization details"
//
// @Success 200 {object} responses.Response{content=organizationDto.Response} "Organization updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [put]
func (oc *OrganizationHandler) Update(c echo.Context) error {
	var request organizationDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := oc.organizationService.Update(organizationUUID, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, organizationMapper.ToResource(updatedOrganization))
}

// Delete an organization
//
// @Summary Delete an organization
// @Description Remove an organization
// @Tags Organizations
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
//
// @Success 204 {object} responses.Response{} "Organization deleted"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [delete]
func (oc *OrganizationHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := oc.organizationService.Delete(organizationUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
