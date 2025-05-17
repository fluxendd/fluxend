package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func UnprocessableResponse(c echo.Context, errors []string) error {
	response := Response{
		Success: false,
		Errors:  errors,
		Content: nil,
	}

	return c.JSON(http.StatusUnprocessableEntity, response)
}
