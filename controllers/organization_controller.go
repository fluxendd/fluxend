package controllers

import (
	"fluxton/pkg"
	"fluxton/requests"
	"fluxton/requests/organization_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type OrganizationController struct {
	organizationService services.OrganizationService
}

func NewOrganizationController(injector *do.Injector) (*OrganizationController, error) {
	organizationService := do.MustInvoke[services.OrganizationService](injector)

	return &OrganizationController{organizationService: organizationService}, nil
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
// @Success 200 {object} responses.Response{content=[]resources.OrganizationResponse} "List of organizations"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations [get]
func (oc *OrganizationController) List(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUserId, _ := pkg.NewAuth(c).Uuid()

	paginationParams := request.ExtractPaginationParams(c)
	organizations, err := oc.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResourceCollection(organizations))
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
// @Success 200 {object} responses.Response{content=resources.OrganizationResponse} "Organization details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [get]
func (oc *OrganizationController) Show(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := pkg.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organization, err := oc.organizationService.GetByID(organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(&organization))
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
// @Success 201 {object} responses.Response{content=resources.OrganizationResponse} "Organization created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations [post]
func (oc *OrganizationController) Store(c echo.Context) error {
	var request organization_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, err := pkg.NewAuth(c).User()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	organization, err := oc.organizationService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.OrganizationResource(&organization))
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
// @Success 200 {object} responses.Response{content=resources.OrganizationResponse} "Organization updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID} [put]
func (oc *OrganizationController) Update(c echo.Context) error {
	var request organization_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := pkg.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := oc.organizationService.Update(organizationUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(updatedOrganization))
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
func (oc *OrganizationController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := pkg.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := oc.organizationService.Delete(organizationUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
