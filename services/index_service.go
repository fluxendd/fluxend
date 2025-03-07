package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type IndexService interface {
	List(tableUUID uuid.UUID, authUser models.AuthUser) ([]string, error)
	GetByName(indexName string, tableUUID uuid.UUID, authUser models.AuthUser) (string, error)
	Create(tableUUID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error)
	Delete(indexName string, tableUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type IndexServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	projectRepo       *repositories.ProjectRepository
	coreTableRepo     *repositories.CoreTableRepository
}

func NewIndexService(injector *do.Injector) (IndexService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &IndexServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
		coreTableRepo:     coreTableRepo,
	}, nil
}

func (s *IndexServiceImpl) List(tableUUID uuid.UUID, authUser models.AuthUser) ([]string, error) {
	table, err := s.coreTableRepo.GetByID(tableUUID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.GetByUUID(table.ProjectUuid)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return nil, errs.NewForbiddenError("project.error.readForbidden")
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return nil, err
	}

	return clientIndexRepo.List(table.Name)
}

func (s *IndexServiceImpl) GetByName(indexName string, tableUUID uuid.UUID, authUser models.AuthUser) (string, error) {
	table, err := s.coreTableRepo.GetByID(tableUUID)
	if err != nil {
		return "", err
	}

	project, err := s.projectRepo.GetByUUID(table.ProjectUuid)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("project.error.readForbidden")
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return "", err
	}

	return clientIndexRepo.GetByName(table.Name, indexName)
}

func (s *IndexServiceImpl) Create(tableUUID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error) {
	table, err := s.coreTableRepo.GetByID(tableUUID)
	if err != nil {
		return "", err
	}

	project, err := s.projectRepo.GetByUUID(table.ProjectUuid)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("table.error.createForbidden")
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return "", err
	}

	hasIndex, err := clientIndexRepo.Has(table.Name, request.Name)
	if err != nil {
		return "", err
	}

	if hasIndex {
		return "", errs.NewUnprocessableError("index.error.alreadyExists")
	}

	_, err = clientIndexRepo.Create(table.Name, request.Name, request.Columns, request.IsUnique)
	if err != nil {
		return "", err
	}

	return clientIndexRepo.GetByName(table.Name, request.Name)
}

func (s *IndexServiceImpl) Delete(indexName string, tableUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	table, err := s.coreTableRepo.GetByID(tableUUID)
	if err != nil {
		return false, err
	}

	project, err := s.projectRepo.GetByUUID(table.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return false, err
	}

	hasIndex, err := clientIndexRepo.Has(table.Name, indexName)
	if err != nil {
		return false, err
	}

	if !hasIndex {
		return false, errs.NewNotFoundError("index.error.notFound")
	}

	return clientIndexRepo.DropIfExists(indexName)
}
