package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeletedResponse(c echo.Context, content interface{}) error {
	response := Response{
		Success: true,
		Errors:  nil,
		Content: content,
	}

	return c.JSON(http.StatusNoContent, response)
}
