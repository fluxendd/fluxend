package handlers

import (
	"fluxend/internal/api/dto"
	userDto "fluxend/internal/api/dto/user"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/config/constants"
	authDomain "fluxend/internal/domain/auth"
	"fluxend/internal/domain/organization"
	"fluxend/internal/domain/user"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type UserHandler struct {
	userService         user.Service
	organizationService organization.Service
}

func NewUserHandler(injector *do.Injector) (*UserHandler, error) {
	userService := do.MustInvoke[user.Service](injector)
	organizationService := do.MustInvoke[organization.Service](injector)

	return &UserHandler{
		userService:         userService,
		organizationService: organizationService,
	}, nil
}

// Show retrieves details of a specific user.
//
// @Summary Retrieve user
// @Description Get details of a specific user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "User UUID"
//
// @Success 200 {object} response.Response{content=user.Response} "User details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /users/{userUUID} [get]
func (uh *UserHandler) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	id, err := request.GetUUIDPathParam(c, "userUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedUser, err := uh.userService.GetByUUID(id)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToUserResource(&fetchedUser))
}

// Me retrieves details of a logged-in user.
//
// @Summary Retrieve logged-in user
// @Description Get details of logged-in specific user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{content=user.Response} "User details"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/me [get]
func (uh *UserHandler) Me(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUserUUID, _ := auth.NewAuth(c).Uuid()
	fetchedUser, err := uh.userService.GetByUUID(authUserUUID)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToUserResource(&fetchedUser))
}

// Login authenticates a user and returns a JWT token.
//
// @Summary Authenticate user
// @Description Authenticate a user and return a JWT token
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param user body user.LoginRequest true "Login request"
//
// @Success 200 {object} response.Response{content=user.Response} "User details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /users/login [post]
func (uh *UserHandler) Login(c echo.Context) error {
	var request userDto.LoginRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	loggedInUser, token, err := uh.userService.Login(userDto.ToLoginUserInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, map[string]interface{}{
		"user":  mapper.ToUserResource(&loggedInUser),
		"token": token,
	})
}

// Store creates a new user.
//
// @Summary Create user
// @Description Add a new user with a name, email, and password
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param user body user.CreateRequest true "User details"
//
// @Success 201 {object} response.Response{content=user.Response} "User created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
//
// @Router /users [post]
func (uh *UserHandler) Store(c echo.Context) error {
	var request userDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	storedUser, token, err := uh.userService.Create(c, userDto.ToCreateUserInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	// Create default organization for the user
	authUser := authDomain.User{
		Uuid:   storedUser.Uuid,
		RoleID: storedUser.RoleID,
	}

	createdOrganization, err := uh.organizationService.Create(constants.DefaultOrganizationName, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, map[string]interface{}{
		"user":  mapper.ToRegisterUserResource(&storedUser, createdOrganization.Uuid),
		"token": token,
	})
}

// Update updates a user.
//
// @Summary Update user
// @Description Update a user's details such as name, email, and password
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param userUUID path string true "User UUID"
// @Param user body user.UpdateRequest true "User details"
//
// @Success 200 {object} response.Response{content=user.Response} "User updated"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /users/{userUUID} [put]
func (uh *UserHandler) Update(c echo.Context) error {
	authUserUUID, err := auth.NewAuth(c).Uuid()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	var request userDto.UpdateRequest
	userUUID, err := request.GetUUIDPathParam(c, "userUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return response.BadRequestResponse(c, "user.error.invalidPayload")
	}

	updatedUser, err := uh.userService.Update(userUUID, authUserUUID, userDto.ToUpdateUserInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToUserResource(updatedUser))
}

// Logout logs out a user by invalidating the JWT token.
//
// @Summary Logout user
// @Description Invalidate the JWT token to log out a user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{} "User logged out"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /users/logout [post]
func (uh *UserHandler) Logout(c echo.Context) error {
	userUUID, err := auth.NewAuth(c).Uuid()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	err = uh.userService.Logout(userUUID)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
