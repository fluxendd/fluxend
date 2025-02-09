package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
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

func (pc *SettingController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	settings, err := pc.settingService.List(authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.SettingResourceCollection(settings))
}

func (pc *SettingController) Update(c echo.Context) error {
	var request requests.SettingUpdateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	updatedSettings, err := pc.settingService.Update(authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.SettingResourceCollection(updatedSettings))
}

func (pc *SettingController) Reset(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	updatedSettings, err := pc.settingService.Reset(authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.SettingResourceCollection(updatedSettings))
}
