package handlers

import (
	"fluxend/internal/api/dto"
	userDto "fluxend/internal/api/dto/user"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/user"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(injector *do.Injector) (*UserHandler, error) {
	userService := do.MustInvoke[user.Service](injector)

	return &UserHandler{userService: userService}, nil
}

// Show retrieves details of a specific user.
//
// @Summary Retrieve
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
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/{userUUID} [get]
func (uh *UserHandler) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	id, err := request.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedUser, err := uh.userService.GetByID(id)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToUserResource(&fetchedUser))
}

// Login authenticates a user and returns a JWT token.
//
// @Summary Authenticate
// @Description Authenticate a user and return a JWT token
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param user body user.LoginRequest true "Login request"
//
// @Success 200 {object} response.Response{content=user.Response} "User details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/login [post]
func (uh *UserHandler) Login(c echo.Context) error {
	var request userDto.LoginRequest
	if err := c.Bind(&request); err != nil {
		return response.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
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
// @Summary Create
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
	if err := c.Bind(&request); err != nil {
		return response.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	storedUser, token, err := uh.userService.Create(c, userDto.ToCreateUserInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, map[string]interface{}{
		"user":  mapper.ToUserResource(&storedUser),
		"token": token,
	})
}

// Update updates a user.
//
// @Summary Update
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
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
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
// @Summary Logout
// @Description Invalidate the JWT token to log out a user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{} "User logged out"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
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
