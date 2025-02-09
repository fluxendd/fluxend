package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/samber/do"
	"time"
)

type SettingService interface {
	List(authUser models.AuthUser) ([]models.Setting, error)
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

func (s *SettingServiceImpl) List(authUser models.AuthUser) ([]models.Setting, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []models.Setting{}, errs.NewForbiddenError("setting.error.listForbidden")
	}

	return s.settingRepo.List()
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
