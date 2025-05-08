package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"time"
)

const settingsCacheKey = "settings"

type SettingService interface {
	List(ctx echo.Context, skipCache bool) ([]models.Setting, error)
	Get(ctx echo.Context, name string) models.Setting
	GetValue(ctx echo.Context, name string) string
	GetBool(ctx echo.Context, name string) bool
	Update(ctx echo.Context, authUser models.AuthUser, request *requests.SettingUpdateRequest) ([]models.Setting, error)
	Reset(ctx echo.Context, authUser models.AuthUser) ([]models.Setting, error)
	GetStorageDriver(ctx echo.Context) string
}

type SettingServiceImpl struct {
	adminPolicy  *policies.AdminPolicy
	databaseRepo *repositories.DatabaseRepository
	settingRepo  *repositories.SettingRepository
}

func NewSettingService(injector *do.Injector) (SettingService, error) {
	policy := policies.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	settingRepo := do.MustInvoke[*repositories.SettingRepository](injector)

	return &SettingServiceImpl{
		adminPolicy:  policy,
		databaseRepo: databaseRepo,
		settingRepo:  settingRepo,
	}, nil
}

func (s *SettingServiceImpl) List(ctx echo.Context, skipCache bool) ([]models.Setting, error) {
	if !skipCache && ctx.Get(settingsCacheKey) != nil {
		settings := ctx.Get(settingsCacheKey).([]models.Setting)

		return settings, nil
	}

	settings, err := s.settingRepo.List()
	if err != nil {
		return nil, err
	}

	ctx.Set(settingsCacheKey, settings)

	return settings, nil
}

func (s *SettingServiceImpl) Get(ctx echo.Context, name string) models.Setting {
	settings, err := s.List(ctx, false)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Error fetching settings")
	}

	for _, setting := range settings {
		if setting.Name == name {
			return setting
		}
	}

	log.Error().
		Str("name", name).
		Msg("Setting not found")

	return models.Setting{}
}

func (s *SettingServiceImpl) GetValue(ctx echo.Context, name string) string {
	setting := s.Get(ctx, name)

	return setting.Value
}

func (s *SettingServiceImpl) GetBool(ctx echo.Context, name string) bool {
	setting := s.Get(ctx, name)

	return setting.Value == "yes"
}

func (s *SettingServiceImpl) Update(ctx echo.Context, authUser models.AuthUser, request *requests.SettingUpdateRequest) ([]models.Setting, error) {
	// Authorization check
	if !s.adminPolicy.CanUpdate(authUser) {
		return nil, errs.NewForbiddenError("setting.error.updateForbidden")
	}

	// Fetch current settings from the repository
	existingSettings, err := s.settingRepo.List()
	if err != nil {
		return nil, err
	}

	// Loop through each setting and update the value
	for _, setting := range request.Settings {
		for i, existingSetting := range existingSettings {
			if existingSetting.Name == setting.Name {
				existingSettings[i].Value = setting.Value
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

func (s *SettingServiceImpl) Reset(ctx echo.Context, authUser models.AuthUser) ([]models.Setting, error) {
	if !s.adminPolicy.CanUpdate(authUser) {
		return []models.Setting{}, errs.NewForbiddenError("setting.error.resetForbidden")
	}

	settings, err := s.settingRepo.List()
	if err != nil {
		return []models.Setting{}, err
	}

	for i := range settings {
		settings[i].Value = settings[i].DefaultValue
		settings[i].UpdatedAt = time.Now()
	}

	_, err = s.settingRepo.Update(settings)
	if err != nil {
		return []models.Setting{}, err
	}

	return s.List(ctx, true)
}

func (s *SettingServiceImpl) GetStorageDriver(ctx echo.Context) string {
	return s.GetValue(ctx, "storageDriver")
}
