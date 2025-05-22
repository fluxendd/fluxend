package logging

import (
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/shared"
	"fluxton/pkg/errors"
	"github.com/samber/do"
)

type Service interface {
	List(paginationParams shared.PaginationParams, authUser auth.User) ([]RequestLog, error)
}

type ServiceImpl struct {
	adminPolicy *admin.Policy
	logRepo     Repository
}

func NewFileService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*admin.Policy](injector)
	logRepo := do.MustInvoke[Repository](injector)

	return &ServiceImpl{
		adminPolicy: policy,
		logRepo:     logRepo,
	}, nil
}

func (s *ServiceImpl) List(paginationParams shared.PaginationParams, authUser auth.User) ([]RequestLog, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []RequestLog{}, errors.NewForbiddenError("log.error.listForbidden")
	}

	return s.logRepo.List(paginationParams)
}
