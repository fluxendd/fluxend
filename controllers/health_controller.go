package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type HealthController struct {
	settingService services.SettingService
}

func NewHealthController(injector *do.Injector) (*HealthController, error) {
	settingService := do.MustInvoke[services.SettingService](injector)

	return &HealthController{settingService: settingService}, nil
}

func (hc *HealthController) Pulse(c echo.Context) error {
	_, err := auth.NewAuth(c).User()
	if err != nil {
		return response.UnauthorizedResponse(c, err.Error())
	}

	// TODO: add logic for health check

	return response.SuccessResponse(c, nil)
}
