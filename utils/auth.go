package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"myapp/models"
)

type Auth struct {
	ctx echo.Context
}

func NewAuth(ctx echo.Context) *Auth {
	return &Auth{ctx: ctx}
}

func (a *Auth) User() (models.AuthenticatedUser, error) {
	user := a.ctx.Get("user")

	// Ensure the value is a map as expected from JWT claims
	userData, ok := user.(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}, errors.New("user.error.invalid_claim_structure")
	}

	return userData, nil
}

func (a *Auth) Id() (uint, error) {
	user, err := a.User()
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (a *Auth) RoleID() (int, error) {
	user, err := a.User()
	if err != nil {
		return 0, err
	}

	return user.RoleID, nil
}
