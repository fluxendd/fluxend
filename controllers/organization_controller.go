package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"myapp/requests"
	"myapp/resources"
	"myapp/responses"
	"myapp/services"
	"myapp/utils"
)

type OrganizationController struct {
	organizationService services.OrganizationService
}

func NewOrganizationController(injector *do.Injector) (*OrganizationController, error) {
	organizationService := do.MustInvoke[services.OrganizationService](injector)

	return &OrganizationController{organizationService: organizationService}, nil
}

func (nc *OrganizationController) List(c echo.Context) error {
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	paginationParams := utils.ExtractPaginationParams(c)
	organizations, err := nc.organizationService.List(paginationParams, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResourceCollection(organizations))
}

func (nc *OrganizationController) Show(c echo.Context) error {
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organization, err := nc.organizationService.GetByID(id, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Store(c echo.Context) error {
	var request requests.OrganizationCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	organization, err := nc.organizationService.Create(&request, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.OrganizationResource(&organization))
}

func (nc *OrganizationController) Update(c echo.Context) error {
	var request requests.OrganizationCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "organization.error.invalidPayload")
	}

	updatedOrganization, err := nc.organizationService.Update(id, authenticatedUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(updatedOrganization))
}

func (nc *OrganizationController) Delete(c echo.Context) error {
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := nc.organizationService.Delete(id, authenticatedUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
