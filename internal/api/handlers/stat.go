package handlers

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/stats"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type StatHandler struct {
	statsService stats.Service
}

func NewStatHandler(injector *do.Injector) (*StatHandler, error) {
	statsService := do.MustInvoke[stats.Service](injector)

	return &StatHandler{statsService: statsService}, nil
}

// Retrieve Retrieves statistics for a project
//
// @Summary Retrieve project statistics
// @Description Get statistics for project
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID query string true "Project UUID"
//
// @Success 200 {object} response.Response{content=stat.Response} "Stats for project"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID}/stats [get]
func (ph *StatHandler) Retrieve(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()
	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedStats, err := ph.statsService.GetAll(projectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToStatResource(&fetchedStats))
}
