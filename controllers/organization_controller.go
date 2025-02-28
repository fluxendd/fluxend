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

func (oc *OrganizationController) List(c echo.Context) error {
	authUserId, _ := utils.NewAuth(c).Uuid()

	paginationParams := utils.ExtractPaginationParams(c)
	organizations, err := oc.organizationService.List(paginationParams, authUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResourceCollection(organizations))
}

func (oc *OrganizationController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organization, err := oc.organizationService.GetByID(organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(&organization))
}

func (oc *OrganizationController) Store(c echo.Context) error {
	var request organization_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, err := utils.NewAuth(c).User()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	organization, err := oc.organizationService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.OrganizationResource(&organization))
}

func (oc *OrganizationController) Update(c echo.Context) error {
	var request organization_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := oc.organizationService.Update(organizationUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.OrganizationResource(updatedOrganization))
}

func (oc *OrganizationController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := oc.organizationService.Delete(organizationUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
