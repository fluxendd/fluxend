package health

import (
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/admin"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/samber/do"
)

type HealthService interface {
	Pulse(authUser models.AuthUser) (Health, error)
}

type HealthServiceImpl struct {
	adminPolicy  *admin.AdminPolicy
	databaseRepo *repositories2.DatabaseRepository
	settingRepo  *repositories2.SettingRepository
}

func NewHealthService(injector *do.Injector) (HealthService, error) {
	policy := admin.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*repositories2.DatabaseRepository](injector)
	settingRepo := do.MustInvoke[*repositories2.SettingRepository](injector)

	return &HealthServiceImpl{
		adminPolicy:  policy,
		databaseRepo: databaseRepo,
		settingRepo:  settingRepo,
	}, nil
}

func (s *HealthServiceImpl) Pulse(authUser models.AuthUser) (Health, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return Health{}, errors.NewForbiddenError("setting.error.listForbidden")
	}

	return Health{}, nil
}
