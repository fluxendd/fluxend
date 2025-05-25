package logging

import (
	"fluxend/internal/domain/admin"
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/shared"
	"fluxend/pkg/errors"
	"github.com/samber/do"
)

type Service interface {
	List(listInput *ListInput, paginationParams shared.PaginationParams, authUser auth.User) ([]RequestLog, error)
}

type ServiceImpl struct {
	adminPolicy *admin.Policy
	logRepo     Repository
}

func NewLogService(injector *do.Injector) (Service, error) {
	policy := admin.NewAdminPolicy()
	logRepo := do.MustInvoke[Repository](injector)

	return &ServiceImpl{
		adminPolicy: policy,
		logRepo:     logRepo,
	}, nil
}

func (s *ServiceImpl) List(listInput *ListInput, paginationParams shared.PaginationParams, authUser auth.User) ([]RequestLog, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []RequestLog{}, errors.NewForbiddenError("log.error.listForbidden")
	}

	return s.logRepo.List(listInput, paginationParams)
}
