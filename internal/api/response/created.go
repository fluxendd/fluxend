package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreatedResponse(c echo.Context, content interface{}) error {
	response := Response{
		Success: true,
		Errors:  nil,
		Content: content,
	}

	return c.JSON(http.StatusCreated, response)
}
