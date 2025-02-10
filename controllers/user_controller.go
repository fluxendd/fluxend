package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
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
	var request requests.UserLoginRequest
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
	var request requests.UserCreateRequest
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
	authUserID := c.Get("userId").(uuid.UUID)
	if authUserID == uuid.Nil {
		return responses.UnauthorizedResponse(c, "user.error.unauthorized")
	}

	var request requests.UserUpdateRequest
	id, err := utils.GetUUIDPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "note.error.invalidPayload")
	}

	updatedUser, err := uc.userService.Update(id, authUserID, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.UserResource(updatedUser))
}

func (uc *UserController) Logout(c echo.Context) error {
	authUserID, err := utils.NewAuth(c).Id()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	err = uc.userService.Logout(authUserID)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
