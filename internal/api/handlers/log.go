package handlers

import (
	loggingDto "fluxend/internal/api/dto/logging"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/logging"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type LogHandler struct {
	logService logging.Service
}

func NewLogHandler(injector *do.Injector) (*LogHandler, error) {
	logService := do.MustInvoke[logging.Service](injector)

	return &LogHandler{logService: logService}, nil
}

// List all logs
//
// @Summary List logs
// @Description Get all logs
// @Tags Logs
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param userUuid query string false "Filter by user UUID"
// @Param status query string false "Filter by status"
// @Param method query string false "Filter by HTTP method"
// @Param endpoint query string false "Filter by endpoint"
// @Param ipAddress query string false "Filter by IP address"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {array} response.Response{content=[]logging.Response} "List of files"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /admin/logs [get]
func (lh *LogHandler) List(c echo.Context) error {
	var request loggingDto.ListRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	logs, err := lh.logService.List(loggingDto.ToLogListInput(&request), paginationParams, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToLoggingResourceCollection(logs))
}
