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

type OrganizationUserController struct {
	organizationService services.OrganizationService
}

func NewOrganizationUserController(injector *do.Injector) (*OrganizationUserController, error) {
	organizationService := do.MustInvoke[services.OrganizationService](injector)

	return &OrganizationUserController{organizationService: organizationService}, nil
}

func (ouc *OrganizationUserController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organizationUsers, err := ouc.organizationService.ListUsers(organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.UserResourceCollection(organizationUsers))
}

func (ouc *OrganizationUserController) Store(c echo.Context) error {
	var request organization_requests.MemberCreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	organizationUser, err := ouc.organizationService.CreateUser(&request, organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.UserResource(&organizationUser))
}

func (ouc *OrganizationUserController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	userID, err := utils.GetUUIDPathParam(c, "userID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := ouc.organizationService.DeleteUser(organizationUUID, userID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
