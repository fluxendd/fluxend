package controllers

import (
	"fluxton/requests/user_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(injector *do.Injector) (*UserController, error) {
	userService := do.MustInvoke[services.UserService](injector)

	return &UserController{userService: userService}, nil
}

func (uc *UserController) Show(c echo.Context) error {
	id, err := utils.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	user, err := uc.userService.GetByID(id)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.UserResource(&user))
}

func (uc *UserController) Login(c echo.Context) error {
	var request user_requests.UserLoginRequest
	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	user, token, err := uc.userService.Login(&request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, map[string]interface{}{
		"user":  resources.UserResource(&user),
		"token": token,
	})
}

func (uc *UserController) Store(c echo.Context) error {
	var request user_requests.UserCreateRequest
	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "user.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	user, err := uc.userService.Create(&request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.UserResource(&user))
}

func (uc *UserController) Update(c echo.Context) error {
	authUserUUID, err := utils.NewAuth(c).Uuid()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	var request user_requests.UserUpdateRequest
	userUUID, err := utils.GetUUIDPathParam(c, "userUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "note.error.invalidPayload")
	}

	updatedUser, err := uc.userService.Update(userUUID, authUserUUID, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.UserResource(updatedUser))
}

func (uc *UserController) Logout(c echo.Context) error {
	userUUID, err := utils.NewAuth(c).Uuid()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	err = uc.userService.Logout(userUUID)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
