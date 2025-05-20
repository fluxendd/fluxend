package handlers

import (
	"fluxton/internal/api/dto/setting"
	settingMapper "fluxton/internal/api/mapper/setting"
	"fluxton/internal/api/response"
	settingDomain "fluxton/internal/domain/setting"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type SettingHandler struct {
	settingService settingDomain.Service
}

func NewSettingHandler(injector *do.Injector) (*SettingHandler, error) {
	settingService := do.MustInvoke[settingDomain.Service](injector)

	return &SettingHandler{settingService: settingService}, nil
}

func (sc *SettingHandler) List(c echo.Context) error {
	settings, err := sc.settingService.List(c, false)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(settings))
}

func (sc *SettingHandler) Update(c echo.Context) error {
	var request setting.UpdateRequest
	authUser, _ := auth.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	updatedSettings, err := sc.settingService.Update(c, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(updatedSettings))
}

func (sc *SettingHandler) Reset(c echo.Context) error {
	authUser, _ := auth.NewAuth(c).User()

	updatedSettings, err := sc.settingService.Reset(c, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(updatedSettings))
}
