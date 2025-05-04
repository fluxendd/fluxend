package responses

import (
	"fluxton/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

func NotFoundResponse(c echo.Context, error string) error {
	response := Response{
		Success: false,
		Errors:  []string{message.Message(error)},
		Content: nil,
	}

	return c.JSON(http.StatusNotFound, response)
}
