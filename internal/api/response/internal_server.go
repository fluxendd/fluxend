package response

import (
	"fluxend/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InternalServerResponse(c echo.Context, error string) error {
	response := InternalServerErrorResponse{
		Success: false,
		Errors:  []string{message.Message(error)},
		Content: nil,
	}

	return c.JSON(http.StatusInternalServerError, response)
}
