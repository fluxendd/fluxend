package responses

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SuccessResponse(c echo.Context, content interface{}) error {
	response := Response{
		Success: true,
		Errors:  nil,
		Content: content,
	}

	return c.JSON(http.StatusOK, response)
}
