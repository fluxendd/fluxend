package handlers

import (
	"fluxton/internal/api/dto"
	organizationDto "fluxton/internal/api/dto/organization"
	userMapper "fluxton/internal/api/mapper/user"
	"fluxton/internal/api/response"
	organizationDomain "fluxton/internal/domain/organization"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type OrganizationMemberHandler struct {
	organizationService organizationDomain.Service
}

func NewOrganizationMemberHandler(injector *do.Injector) (*OrganizationMemberHandler, error) {
	organizationService := do.MustInvoke[organizationDomain.Service](injector)

	return &OrganizationMemberHandler{organizationService: organizationService}, nil
}

// List all users in an organization
//
// @Summary List all users in an organization
// @Description Get all users in an organization
// @Tags Organization Members
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
//
// @Success 201 {object} responses.Response{content=[]resources.UserResponse} "User created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID}/users [get]
func (ouc *OrganizationMemberHandler) List(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	organizationUsers, err := ouc.organizationService.ListUsers(organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, userMapper.ToResponseCollection(organizationUsers))
}

// Store creates a user in an organization
//
// @Summary Create a user in an organization
// @Description Add a new user to an organization
// @Tags Organization Members
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
// @Param user body organization_requests.MemberCreateRequest true "User ID JSON"
//
// @Success 201 {object} responses.Response{content=resources.UserResponse} "User created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID}/users [post]
func (ouc *OrganizationMemberHandler) Store(c echo.Context) error {
	var request organizationDto.MemberCreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	organizationUser, err := ouc.organizationService.CreateUser(request.UserID, organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, userMapper.ToResponse(&organizationUser))
}

// Delete a user from an organization
//
// @Summary Delete a user from an organization
// @Description Remove a user from an organization
// @Tags Organization Members
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organization_id path string true "Organization ID"
// @Param user_id path string true "User ID"
//
// @Success 204 {object} nil "User deleted"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /organizations/{organizationUUID}/users/{userUUID} [delete]
func (ouc *OrganizationMemberHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDPathParam(c, "organizationUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	userID, err := request.GetUUIDPathParam(c, "userID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if err := ouc.organizationService.DeleteUser(organizationUUID, userID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
