package controllers

import (
	"fluxton/requests/organization_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
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

func (nc *OrganizationController) List(c echo.Context) error {
	authUserId, _ := utils.NewAuth(c).Uuid()

	paginationParams := utils.ExtractPaginationParams(c)
	organizations, err := nc.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResourceCollection(organizations))
}

func (nc *OrganizationController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organization, err := nc.organizationService.GetByID(organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Store(c echo.Context) error {
	var request organization_requests.OrganizationCreateRequest
	authUser, err := utils.NewAuth(c).User()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	organization, err := nc.organizationService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Update(c echo.Context) error {
	var request organization_requests.OrganizationCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	updatedOrganization, err := nc.organizationService.Update(organizationUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(updatedOrganization))
}

func (nc *OrganizationController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := nc.organizationService.Delete(organizationUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
