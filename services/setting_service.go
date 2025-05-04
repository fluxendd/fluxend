package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"time"
)

type SettingService interface {
	List() ([]models.Setting, error)
	Get(name string) models.Setting
	GetValue(name string) string
	GetBool(name string) bool
	Update(authUser models.AuthUser, request *requests.SettingUpdateRequest) ([]models.Setting, error)
	Reset(authUser models.AuthUser) ([]models.Setting, error)
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

func (s *SettingServiceImpl) List() ([]models.Setting, error) {
	return s.settingRepo.List()
}

func (s *SettingServiceImpl) Get(name string) models.Setting {
	settings, err := s.settingRepo.List() // cached settings so no need to worry about multiple calls
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

func (s *SettingServiceImpl) GetValue(name string) string {
	setting := s.Get(name)

	return setting.Value
}

func (s *SettingServiceImpl) GetBool(name string) bool {
	setting := s.Get(name)

	return setting.Value == "yes"
}

func (s *SettingServiceImpl) Update(authUser models.AuthUser, request *requests.SettingUpdateRequest) ([]models.Setting, error) {
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

	return s.settingRepo.List()
}

func (s *SettingServiceImpl) Reset(authUser models.AuthUser) ([]models.Setting, error) {
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

	return s.settingRepo.List()
}
