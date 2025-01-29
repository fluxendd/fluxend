package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	ctx echo.Context
}

func NewAuth(ctx echo.Context) *Auth {
	return &Auth{ctx: ctx}
}

func (a *Auth) Id() (uint, error) {
	authenticatedUserId := uint(a.ctx.Get("id").(float64))
	if authenticatedUserId == 0 {
		return 0, errors.New("user.error.unauthenticated")
	}

	return authenticatedUserId, nil
}

func (a *Auth) Email() (string, error) {
	authenticatedUserEmail := a.ctx.Get("email").(string)
	if authenticatedUserEmail == "" {
		return "", errors.New("user.error.unauthenticated")
	}

	return authenticatedUserEmail, nil
}
