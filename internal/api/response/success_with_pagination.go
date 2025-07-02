package response

import (
	"fluxend/internal/domain/shared"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SuccessResponseWithPagination(c echo.Context, content interface{}, paginationDetails shared.PaginationDetails) error {
	response := Response{
		Success: true,
		Errors:  nil,
		Content: content,
		Metadata: map[string]interface{}{
			"total": paginationDetails.Total,
			"page":  paginationDetails.Page,
			"limit": paginationDetails.Limit,
		},
	}

	return c.JSON(http.StatusOK, response)
}
