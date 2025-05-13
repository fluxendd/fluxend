package email

import (
	"github.com/labstack/echo/v4"
)

type EmailInterface interface {
	Send(ctx echo.Context, to, subject, body string) error
}
