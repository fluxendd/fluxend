package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"github.com/samber/do"
)

type HealthService interface {
	Pulse(authUser models.AuthUser) (models.Health, error)
}

type HealthServiceImpl struct {
	adminPolicy  *policies.AdminPolicy
	databaseRepo *repositories.DatabaseRepository
	settingRepo  *repositories.SettingRepository
}

func NewHealthService(injector *do.Injector) (HealthService, error) {
	policy := policies.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	settingRepo := do.MustInvoke[*repositories.SettingRepository](injector)

	return &HealthServiceImpl{
		adminPolicy:  policy,
		databaseRepo: databaseRepo,
		settingRepo:  settingRepo,
	}, nil
}

func (s *HealthServiceImpl) Pulse(authUser models.AuthUser) (models.Health, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return models.Health{}, errs.NewForbiddenError("setting.error.listForbidden")
	}

	return models.Health{}, nil
}
