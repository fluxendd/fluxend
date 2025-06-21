package response

import (
	"fluxend/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

func BadRequestResponse(c echo.Context, error string) error {
	response := BadRequestErrorResponse{
		Success: false,
		Errors:  []string{message.Message(error)},
		Content: nil,
	}

	return c.JSON(http.StatusBadRequest, response)
}
