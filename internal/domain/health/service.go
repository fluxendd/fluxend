package health

import (
	"fluxend/internal/domain/admin"
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/setting"
	"fluxend/internal/domain/shared"
	"fluxend/pkg/errors"
	"github.com/samber/do"
)

const statusOk = "OK"
const statusError = "ERROR"

type Service interface {
	Pulse(authUser auth.User) (Health, error)
}

type ServiceImpl struct {
	adminPolicy      *admin.Policy
	settingRepo      setting.Repository
	projectRepo      project.Repository
	postgrestService shared.PostgrestService
}

func NewHealthService(injector *do.Injector) (Service, error) {
	policy := admin.NewAdminPolicy()
	settingRepo := do.MustInvoke[setting.Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)
	postgrestService := do.MustInvoke[shared.PostgrestService](injector)

	return &ServiceImpl{
		adminPolicy:      policy,
		settingRepo:      settingRepo,
		projectRepo:      projectRepo,
		postgrestService: postgrestService,
	}, nil
}

func (s *ServiceImpl) Pulse(authUser auth.User) (Health, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return Health{}, errors.NewForbiddenError("setting.error.listForbidden")
	}

	response := Health{
		DatabaseStatus:  statusOk,
		AppStatus:       statusOk,
		PostgrestStatus: statusOk,
	}

	allProjects, err := s.projectRepo.List(shared.PaginationParams{})
	if err != nil {
		return Health{}, err
	}

	for _, currentProject := range allProjects {
		if !s.postgrestService.HasContainer(currentProject.DBName) {
			response.PostgrestStatus = statusError

			break
		}
	}

	return response, nil
}
