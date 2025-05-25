package handlers

import (
	settingDto "fluxend/internal/api/dto/setting"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/setting"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type SettingHandler struct {
	settingService setting.Service
}

func NewSettingHandler(injector *do.Injector) (*SettingHandler, error) {
	settingService := do.MustInvoke[setting.Service](injector)

	return &SettingHandler{settingService: settingService}, nil
}

func (sh *SettingHandler) List(c echo.Context) error {
	settings, err := sh.settingService.List()
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToSettingResourceCollection(settings))
}

func (sh *SettingHandler) Update(c echo.Context) error {
	var request settingDto.UpdateRequest
	authUser, _ := auth.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	updatedSettings, err := sh.settingService.Update(authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToSettingResourceCollection(updatedSettings))
}

func (sh *SettingHandler) Reset(c echo.Context) error {
	authUser, _ := auth.NewAuth(c).User()

	updatedSettings, err := sh.settingService.Reset(authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToSettingResourceCollection(updatedSettings))
}
