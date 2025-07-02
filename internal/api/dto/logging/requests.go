package logging

import (
	"fluxend/internal/api/dto"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/labstack/echo/v4"
	"time"
)

type ListRequest struct {
	dto.BaseRequest
	UserUuid  uuid.NullUUID `query:"userUuid"`
	Status    null.String   `query:"status"`
	Method    null.String   `query:"method"`
	Endpoint  null.String   `query:"endpoint"`
	IPAddress null.String   `query:"ipAddress"`
	DateStart null.String   `query:"dateStart"`
	DateEnd   null.String   `query:"dateEnd"`

	Limit int    `query:"limit"`
	Page  int    `query:"page"`
	Sort  string `query:"sort"`
	Order string `query:"order"`

	// Add parsed time fields for internal use
	ParsedDateStart *time.Time
	ParsedDateEnd   *time.Time
}

func (r *ListRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	var dateStart, dateEnd time.Time
	var hasDateStart, hasDateEnd bool

	if r.DateStart.Valid {
		parsed, err := time.Parse("02-01-2006", r.DateStart.String)
		if err != nil {
			return []string{"Invalid dateStart format, expected DD-MM-YYYY"}
		}

		dateStart = parsed
		r.ParsedDateStart = &dateStart
		hasDateStart = true
	}

	if r.DateEnd.Valid {
		parsed, err := time.Parse("02-01-2006", r.DateEnd.String)
		if err != nil {
			return []string{"Invalid dateEnd format, expected DD-MM-YYYY"}
		}

		dateEnd = parsed
		r.ParsedDateEnd = &dateEnd
		hasDateEnd = true
	}

	// Check if dateEnd is after dateStart (only if both are valid)
	if hasDateStart && hasDateEnd {
		if !dateEnd.After(dateStart) && !dateEnd.Equal(dateStart) {
			return []string{"dateEnd must be after dateStart"}
		}
	}

	return nil
}
