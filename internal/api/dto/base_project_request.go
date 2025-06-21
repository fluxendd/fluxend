package dto

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DefaultRequestWithProjectHeader struct {
	BaseRequest
	ProjectUUID uuid.UUID `json:"projectUUID"`
}

func (r *DefaultRequestWithProjectHeader) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	if err := r.WithProjectHeader(c); err != nil {
		return []string{err.Error()}
	}

	return nil
}

func (r *DefaultRequestWithProjectHeader) WithProjectHeader(c echo.Context) error {
	projectUUID, err := uuid.Parse(c.Request().Header.Get("X-Project"))
	if err != nil {
		return errors.New("invalid project UUID")
	}

	r.ProjectUUID = projectUUID

	return nil
}
