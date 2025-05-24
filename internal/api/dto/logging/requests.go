package logging

import (
	"fluxton/internal/api/dto"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/labstack/echo/v4"
)

type ListRequest struct {
	dto.BaseRequest
	UserUuid  uuid.NullUUID `query:"userUuid"`
	Status    null.String   `query:"status"`
	Method    null.String   `query:"method"`
	Endpoint  null.String   `query:"endpoint"`
	IPAddress null.String   `query:"ipAddress"`

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
