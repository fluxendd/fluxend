package response

import (
	"errors"
	flxErrors "fluxend/pkg/errors"
	"github.com/labstack/echo/v4"
)

func ErrorResponse(c echo.Context, err error) error {
	var notFoundErr *flxErrors.NotFoundError
	var unauthorizedErr *flxErrors.UnauthorizedError
	var forbiddenErr *flxErrors.ForbiddenError
	var badRequestErr *flxErrors.BadRequestError

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
