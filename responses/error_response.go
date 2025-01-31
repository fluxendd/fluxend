package responses

import (
	"errors"
	"fluxton/errs"
	"github.com/labstack/echo/v4"
)

func ErrorResponse(c echo.Context, err error) error {
	var notFoundErr *errs.NotFoundError
	var unauthorizedErr *errs.UnauthorizedError
	var forbiddenErr *errs.ForbiddenError
	var badRequestErr *errs.BadRequestError

	if errors.As(err, &notFoundErr) {
		return NotFoundResponse(c, err.Error())
	}

	if errors.As(err, &unauthorizedErr) {
		return UnauthorizedResponse(c, err.Error())
	}

	if errors.As(err, &forbiddenErr) {
		return ForbiddenResponse(c, err.Error())
	}

	if errors.As(err, &badRequestErr) {
		return BadRequestResponse(c, err.Error())
	}

	return InternalServerResponse(c, err.Error())
}
