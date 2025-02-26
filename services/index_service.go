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
	List(tableID, projectID uuid.UUID, authUser models.AuthUser) ([]string, error)
	GetByName(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (string, error)
	Create(projectID, tableID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error)
	Delete(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error)
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

func (s *IndexServiceImpl) List(tableID, projectID uuid.UUID, authUser models.AuthUser) ([]string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return nil, errs.NewForbiddenError("project.error.readForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return nil, err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return nil, err
	}

	return clientIndexRepo.List(table.Name)
}

func (s *IndexServiceImpl) GetByName(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("project.error.readForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return "", err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return "", err
	}

	return clientIndexRepo.GetByName(table.Name, indexName)
}

func (s *IndexServiceImpl) Create(projectID, tableID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("table.error.createForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return "", err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return "", err
	}

	_, err = clientIndexRepo.Create(table.Name, request.Name, request.Columns, request.IsUnique)
	if err != nil {
		return "", err
	}

	return clientIndexRepo.GetByName(table.Name, request.Name)
}

func (s *IndexServiceImpl) Delete(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	_, err = s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return false, err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return false, err
	}

	return clientIndexRepo.DropIfExists(indexName)
}
