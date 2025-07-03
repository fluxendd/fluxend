package handlers

import (
	loggingDto "fluxend/internal/api/dto/logging"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/logging"
	"fluxend/pkg/auth"
	"github.com/google/uuid"
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
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Param userUuid query string false "Filter by user UUID"
// @Param status query string false "Filter by status"
// @Param method query string false "Filter by HTTP method"
// @Param endpoint query string false "Filter by endpoint"
// @Param ipAddress query string false "Filter by IP address"
// @Param startTime query string false "Filter after a unix timestamp"
// @Param endTime query string false "Filter before a unix timestamp"
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} response.Response{content=[]logging.Response} "List of logs"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects/{projectUUID}/logs [get]
func (lh *LogHandler) List(c echo.Context) error {
	var request loggingDto.ListRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	paginationParams := request.ExtractPaginationParams(c)
	input := loggingDto.ToLogListInput(&request, uuid.NullUUID{Valid: true, UUID: projectUUID})

	logs, paginationDetails, err := lh.logService.List(input, paginationParams, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponseWithPagination(c, mapper.ToLoggingResourceCollection(logs), paginationDetails)
}
