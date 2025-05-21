package setting

import (
	"fluxton/internal/api/dto/setting"
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"time"
)

const settingsCacheKey = "settings"

type Service interface {
	List() ([]Setting, error)
	Get(name string) Setting
	GetValue(name string) string
	GetBool(name string) bool
	Update(authUser auth.User, request *setting.UpdateRequest) ([]Setting, error)
	Reset(authUser auth.User) ([]Setting, error)
	GetStorageDriver() string
}

type ServiceImpl struct {
	adminPolicy *admin.Policy
	settingRepo Repository
}

func NewSettingService(injector *do.Injector) (Service, error) {
	policy := admin.NewAdminPolicy()
	settingRepo := do.MustInvoke[Repository](injector)

	return &ServiceImpl{
		adminPolicy: policy,
		settingRepo: settingRepo,
	}, nil
}

func (s *ServiceImpl) List() ([]Setting, error) {
	// TODO: cache using Redis instead of context
	settings, err := s.settingRepo.List()
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *ServiceImpl) Get(name string) Setting {
	fetchedSetting, err := s.settingRepo.Get(name)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Error fetching setting")

		return Setting{}
	}

	return fetchedSetting
}

func (s *ServiceImpl) GetValue(name string) string {
	currentSetting := s.Get(name)

	return currentSetting.Value
}

func (s *ServiceImpl) GetBool(name string) bool {
	currentSetting := s.Get(name)

	return currentSetting.Value == "yes"
}

func (s *ServiceImpl) Update(authUser auth.User, request *setting.UpdateRequest) ([]Setting, error) {
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

	return s.List()
}

func (s *ServiceImpl) Reset(authUser auth.User) ([]Setting, error) {
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

	return s.List()
}

func (s *ServiceImpl) GetStorageDriver() string {
	return s.GetValue("storageDriver")
}
