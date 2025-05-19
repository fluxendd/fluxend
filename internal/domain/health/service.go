package health

import (
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/database"
	"fluxton/internal/domain/setting"
	"fluxton/pkg/errors"
	"github.com/samber/do"
)

type Service interface {
	Pulse(authUser auth.User) (Health, error)
}

type ServiceImpl struct {
	adminPolicy  *admin.Policy
	databaseRepo *database.DatabaseService
	settingRepo  *setting.Repository
}

func NewHealthService(injector *do.Injector) (Service, error) {
	policy := admin.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*database.DatabaseService](injector)
	settingRepo := do.MustInvoke[*setting.Repository](injector)

	return &ServiceImpl{
		adminPolicy:  policy,
		databaseRepo: databaseRepo,
		settingRepo:  settingRepo,
	}, nil
}

func (s *ServiceImpl) Pulse(authUser auth.User) (Health, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return Health{}, errors.NewForbiddenError("setting.error.listForbidden")
	}

	return Health{}, nil
}
