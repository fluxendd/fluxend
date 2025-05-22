package handlers

import (
	"fluxton/internal/api/dto/logging"
	logMapper "fluxton/internal/api/mapper/logging"
	"fluxton/internal/api/response"
	logDomain "fluxton/internal/domain/logging"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type LogHandler struct {
	logService logDomain.Service
}

func NewLogHandler(injector *do.Injector) (*LogHandler, error) {
	logService := do.MustInvoke[logDomain.Service](injector)

	return &LogHandler{logService: logService}, nil
}

// List all logs
//
// @Summary List all logs
// @Description Get all logs
// @Tags Logs
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {array} response.Response{content=[]logging.Response} "List of files"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /admin/logs [get]
func (lh *LogHandler) List(c echo.Context) error {
	var request logging.ListRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	logs, err := lh.logService.List(paginationParams, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, logMapper.ToResourceCollection(logs))
}
