package response

import (
	"fluxend/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ForbiddenResponse(c echo.Context, error string) error {
	response := ForbiddenErrorResponse{
		Success: false,
		Errors:  []string{message.Message(error)},
		Content: nil,
	}

	return c.JSON(http.StatusForbidden, response)
}
