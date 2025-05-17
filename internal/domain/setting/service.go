package setting

import (
	"fluxton/internal/api/dto/setting"
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"time"
)

const settingsCacheKey = "settings"

type SettingService interface {
	List(ctx echo.Context, skipCache bool) ([]Setting, error)
	Get(ctx echo.Context, name string) Setting
	GetValue(ctx echo.Context, name string) string
	GetBool(ctx echo.Context, name string) bool
	Update(ctx echo.Context, authUser auth.AuthUser, request *setting.UpdateRequest) ([]Setting, error)
	Reset(ctx echo.Context, authUser auth.AuthUser) ([]Setting, error)
	GetStorageDriver(ctx echo.Context) string
}

type SettingServiceImpl struct {
	adminPolicy *admin.AdminPolicy
	settingRepo *Repository
}

func NewSettingService(injector *do.Injector) (SettingService, error) {
	policy := admin.NewAdminPolicy()
	settingRepo := do.MustInvoke[*Repository](injector)

	return &SettingServiceImpl{
		adminPolicy: policy,
		settingRepo: settingRepo,
	}, nil
}

func (s *SettingServiceImpl) List(ctx echo.Context, skipCache bool) ([]Setting, error) {
	if !skipCache && ctx.Get(settingsCacheKey) != nil {
		settings := ctx.Get(settingsCacheKey).([]Setting)

		return settings, nil
	}

	settings, err := s.settingRepo.List()
	if err != nil {
		return nil, err
	}

	ctx.Set(settingsCacheKey, settings)

	return settings, nil
}

func (s *SettingServiceImpl) Get(ctx echo.Context, name string) Setting {
	settings, err := s.List(ctx, false)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Error fetching settings")
	}

	for _, currentSetting := range settings {
		if currentSetting.Name == name {
			return currentSetting
		}
	}

	log.Error().
		Str("name", name).
		Msg("Setting not found")

	return Setting{}
}

func (s *SettingServiceImpl) GetValue(ctx echo.Context, name string) string {
	currentSetting := s.Get(ctx, name)

	return currentSetting.Value
}

func (s *SettingServiceImpl) GetBool(ctx echo.Context, name string) bool {
	currentSetting := s.Get(ctx, name)

	return currentSetting.Value == "yes"
}

func (s *SettingServiceImpl) Update(ctx echo.Context, authUser auth.AuthUser, request *setting.UpdateRequest) ([]Setting, error) {
	// Authorization check
	if !s.adminPolicy.CanUpdate(authUser) {
		return nil, errors.NewForbiddenError("setting.error.updateForbidden")
	}

	// Fetch current settings from the repository
	existingSettings, err := s.settingRepo.List()
	if err != nil {
		return nil, err
	}

	// Loop through each setting and update the value
	for _, currentSetting := range request.Settings {
		for i, existingSetting := range existingSettings {
			if existingSetting.Name == currentSetting.Name {
				existingSettings[i].Value = currentSetting.Value
				existingSettings[i].UpdatedAt = time.Now()
			}
		}
	}

	// Update the settings in the repository
	_, err = s.settingRepo.Update(existingSettings)
	if err != nil {
		return nil, err
	}

	return s.List(ctx, true)
}

func (s *SettingServiceImpl) Reset(ctx echo.Context, authUser auth.AuthUser) ([]Setting, error) {
	if !s.adminPolicy.CanUpdate(authUser) {
		return []Setting{}, errors.NewForbiddenError("setting.error.resetForbidden")
	}

	settings, err := s.settingRepo.List()
	if err != nil {
		return []Setting{}, err
	}

	for i := range settings {
		settings[i].Value = settings[i].DefaultValue
		settings[i].UpdatedAt = time.Now()
	}

	_, err = s.settingRepo.Update(settings)
	if err != nil {
		return []Setting{}, err
	}

	return s.List(ctx, true)
}

func (s *SettingServiceImpl) GetStorageDriver(ctx echo.Context) string {
	return s.GetValue(ctx, "storageDriver")
}
