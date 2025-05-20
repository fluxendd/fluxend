package handlers

import (
	"fluxton/internal/api/response"
	"fluxton/internal/domain/setting"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type HealthHandler struct {
	settingService setting.Service
}

func NewHealthHandler(injector *do.Injector) (*HealthHandler, error) {
	settingService := do.MustInvoke[setting.Service](injector)

	return &HealthHandler{settingService: settingService}, nil
}

func (hc *HealthHandler) Pulse(c echo.Context) error {
	_, err := auth.NewAuth(c).User()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	// TODO: add logic for health check

	return response.SuccessResponse(c, nil)
}
