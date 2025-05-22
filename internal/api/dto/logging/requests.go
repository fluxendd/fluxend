package logging

import (
	"fluxton/internal/api/dto"
	"github.com/labstack/echo/v4"
)

type ListRequest struct {
	dto.BaseRequest
	UserUuid  string `query:"userUuid"`
	Status    string `query:"status"`
	Method    string `query:"method"`
	Endpoint  string `query:"endpoint"`
	IPAddress string `query:"ipAddress"`

	Limit int    `query:"limit"`
	Page  int    `query:"page"`
	Sort  string `query:"sort"`
	Order string `query:"order"`
}

func (r *ListRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	return nil
}
