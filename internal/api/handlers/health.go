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
