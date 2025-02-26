package utils

import (
	"errors"
	"fluxton/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	ctx echo.Context
}

func NewAuth(ctx echo.Context) *Auth {
	return &Auth{ctx: ctx}
}

func (a *Auth) User() (models.AuthUser, error) {
	user := a.ctx.Get("user")

	// Ensure the value is a map as expected from JWT claims
	userData, ok := user.(models.AuthUser)
	if !ok {
		return models.AuthUser{}, errors.New("user.error.invalid_claim_structure")
	}

	return userData, nil
}

func (a *Auth) Uuid() (uuid.UUID, error) {
	user, err := a.User()
	if err != nil {
		return uuid.Nil, err
	}

	return user.Uuid, nil
}

func (a *Auth) RoleID() (int, error) {
	user, err := a.User()
	if err != nil {
		return 0, err
	}

	return user.RoleID, nil
}
