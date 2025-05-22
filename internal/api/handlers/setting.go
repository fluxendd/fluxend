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

func (sh *SettingHandler) List(c echo.Context) error {
	settings, err := sh.settingService.List()
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(settings))
}

func (sh *SettingHandler) Update(c echo.Context) error {
	var request setting.UpdateRequest
	authUser, _ := auth.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	updatedSettings, err := sh.settingService.Update(authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(updatedSettings))
}

func (sh *SettingHandler) Reset(c echo.Context) error {
	authUser, _ := auth.NewAuth(c).User()

	updatedSettings, err := sh.settingService.Reset(authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, settingMapper.ToResourceCollection(updatedSettings))
}
