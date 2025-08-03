package logging

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	"fluxend/pkg/errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

type Service interface {
	List(
		listInput *ListInput,
		paginationParams shared.PaginationParams,
		authUser auth.User,
	) ([]RequestLog, shared.PaginationDetails, error)
	Store(storeInput *StoreInput)
}

type ServiceImpl struct {
	projectPolicy *project.Policy
	logRepo       Repository
	projectRepo   project.Repository
}

func NewLogService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	logRepo := do.MustInvoke[Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		projectPolicy: policy,
		logRepo:       logRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *ServiceImpl) List(listInput *ListInput, paginationParams shared.PaginationParams, authUser auth.User) ([]RequestLog, shared.PaginationDetails, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(listInput.ProjectUuid.UUID)
	if err != nil {
		return nil, shared.PaginationDetails{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return nil, shared.PaginationDetails{}, errors.NewForbiddenError("project.error.viewForbidden")
	}

	return s.logRepo.List(listInput, paginationParams)
}

func (s *ServiceImpl) Store(storeInput *StoreInput) {
	pkg.DumpJSON(storeInput)
	if storeInput == nil {
		fmt.Println("Store input is nil")
		return
	}

	if storeInput.IPAddress == "" {
		storeInput.IPAddress = "0.0.0.0"
	}

	projectUUID, err := s.projectRepo.GetUUIDByDatabaseName(storeInput.DbName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch project by database name")

		return
	}

	storeInput.ProjectUUID = projectUUID

	err = s.logRepo.CreatePostgrest(storeInput)
	if err != nil {
		log.Error().Err(err).Msg("Failed to store log")
	}
}
