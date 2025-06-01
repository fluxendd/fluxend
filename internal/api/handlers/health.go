package handlers

import (
	"fluxend/internal/api/response"
	"fluxend/internal/domain/health"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type HealthHandler struct {
	healthService health.Service
}

func NewHealthHandler(injector *do.Injector) (*HealthHandler, error) {
	healthService := do.MustInvoke[health.Service](injector)

	return &HealthHandler{healthService: healthService}, nil
}

// Pulse Health check endpoint
//
// @Summary Check system health
// @Description Check the health status of the system
// @Tags Admin
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{content=dto.GenericResponse} "Health status"
// @Failure 401 "Unauthorized"
//
// @Router /admin/health [get]
func (hh *HealthHandler) Pulse(c echo.Context) error {
	authUser, err := auth.NewAuth(c).User()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	status, err := hh.healthService.Pulse(authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, status)
}
