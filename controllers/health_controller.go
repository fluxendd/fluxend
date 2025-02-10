package controllers

import (
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
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

func (pc *HealthController) Pulse(c echo.Context) error {
	_, err := utils.NewAuth(c).User()
	if err != nil {
		return responses.UnauthorizedResponse(c, err.Error())
	}

	// TODO: add logic for health check

	return responses.SuccessResponse(c, nil)
}
