package handlers

import (
	"fluxton/internal/api/dto"
	userDto "fluxton/internal/api/dto/user"
	userMapper "fluxton/internal/api/mapper/user"
	"fluxton/internal/api/response"
	userDomain "fluxton/internal/domain/user"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type UserController struct {
	userService userDomain.UserService
}

func NewUserHandler(injector *do.Injector) (*UserController, error) {
	userService := do.MustInvoke[userDomain.UserService](injector)

	return &UserController{userService: userService}, nil
}

// Show retrieves details of a specific user.
//
// @Summary Show details of a single user
// @Description Get details of a specific user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param id path string true "User UUID"
//
// @Success 200 {object} responses.Response{content=resources.UserResponse} "User details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/{userUUID} [get]
func (uc *UserController) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	id, err := request.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	user, err := uc.userService.GetByID(id)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, userMapper.ToResponse(&user))
}

// Login authenticates a user and returns a JWT token.
//
// @Summary Authenticate a user
// @Description Authenticate a user and return a JWT token
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param user body user.LoginRequest true "Login request"
//
// @Success 200 {object} responses.Response{content=resources.UserResponse} "User details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/login [post]
func (uc *UserController) Login(c echo.Context) error {
	var request userDto.LoginRequest
	if err := c.Bind(&request); err != nil {
		return response.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	user, token, err := uc.userService.Login(&request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, map[string]interface{}{
		"user":  userMapper.ToResponse(&user),
		"token": token,
	})
}

// Store creates a new user.
//
// @Summary Create a new user
// @Description Add a new user with a name, email, and password
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param user body user.CreateRequest true "User details"
//
// @Success 201 {object} responses.Response{content=resources.UserResponse} "User created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 500 "Internal server error"
//
// @Router /users [post]
func (uc *UserController) Store(c echo.Context) error {
	var request userDto.CreateRequest
	if err := c.Bind(&request); err != nil {
		return response.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	user, token, err := uc.userService.Create(c, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, map[string]interface{}{
		"user":  userMapper.ToResponse(&user),
		"token": token,
	})
}

// Update updates a user.
//
// @Summary Update a user
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
// @Success 200 {object} responses.Response{content=resources.UserResponse} "User updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/{userUUID} [put]
func (uc *UserController) Update(c echo.Context) error {
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

	updatedUser, err := uc.userService.Update(userUUID, authUserUUID, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, userMapper.ToResponse(updatedUser))
}

// Logout logs out a user by invalidating the JWT token.
//
// @Summary Logout a user
// @Description Invalidate the JWT token to log out a user
// @Tags Users
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} responses.Response{} "User logged out"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /users/logout [post]
func (uc *UserController) Logout(c echo.Context) error {
	userUUID, err := auth.NewAuth(c).Uuid()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	err = uc.userService.Logout(userUUID)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
