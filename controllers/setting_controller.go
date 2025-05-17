package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type SettingController struct {
	settingService services.SettingService
}

func NewSettingController(injector *do.Injector) (*SettingController, error) {
	settingService := do.MustInvoke[services.SettingService](injector)

	return &SettingController{settingService: settingService}, nil
}

func (sc *SettingController) List(c echo.Context) error {
	settings, err := sc.settingService.List(c, false)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.SettingResourceCollection(settings))
}

func (sc *SettingController) Update(c echo.Context) error {
	var request requests.SettingUpdateRequest
	authUser, _ := auth.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	updatedSettings, err := sc.settingService.Update(c, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.SettingResourceCollection(updatedSettings))
}

func (sc *SettingController) Reset(c echo.Context) error {
	authUser, _ := auth.NewAuth(c).User()

	updatedSettings, err := sc.settingService.Reset(c, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.SettingResourceCollection(updatedSettings))
}
