package response

import (
	"fluxton/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InternalServerResponse(c echo.Context, error string) error {
	response := Response{
		Success: false,
		Errors:  []string{message.Message(error)},
		Content: nil,
	}

	return c.JSON(http.StatusInternalServerError, response)
}
