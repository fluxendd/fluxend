package health

import (
	"fluxton/internal/adapters/client"
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/setting"
	"fluxton/pkg/errors"
	"github.com/samber/do"
)

type HealthService interface {
	Pulse(authUser auth.User) (Health, error)
}

type HealthServiceImpl struct {
	adminPolicy  *admin.Policy
	databaseRepo *client.Repository
	settingRepo  *setting.Repository
}

func NewHealthService(injector *do.Injector) (HealthService, error) {
	policy := admin.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*client.Repository](injector)
	settingRepo := do.MustInvoke[*setting.Repository](injector)

	return &HealthServiceImpl{
		adminPolicy:  policy,
		databaseRepo: databaseRepo,
		settingRepo:  settingRepo,
	}, nil
}

func (s *HealthServiceImpl) Pulse(authUser auth.User) (Health, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return Health{}, errors.NewForbiddenError("setting.error.listForbidden")
	}

	return Health{}, nil
}
