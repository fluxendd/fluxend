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

// List Settings
//
// @Summary List settings
// @Description Retrieve all settings for the application
// @Tags Admin
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{content=[]setting.Response} "List of indexes"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /admin/settings [get]
func (sh *SettingHandler) List(c echo.Context) error {
	settings, err := sh.settingService.List()
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToSettingResourceCollection(settings))
}

// Update Settings
//
// @Summary Update settings
// @Description Update application settings
// @Tags Admin
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Param form body setting.UpdateRequest true "Settings update request"
//
// @Success 200 {object} response.Response{content=[]setting.Response} "Form updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /admin/settings [put]
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

// Reset Settings
//
// @Summary Reset settings
// @Description Reset application settings
// @Tags Admin
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
//
// @Success 200 {object} response.Response{content=[]setting.Response} "Settings reset"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /admin/settings/reset [put]
func (sh *SettingHandler) Reset(c echo.Context) error {
	authUser, _ := auth.NewAuth(c).User()

	updatedSettings, err := sh.settingService.Reset(authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToSettingResourceCollection(updatedSettings))
}
