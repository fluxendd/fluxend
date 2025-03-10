package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type IndexService interface {
	List(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) ([]string, error)
	GetByName(indexName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (string, error)
	Create(fullTableName string, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error)
	Delete(indexName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type IndexServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	projectRepo       *repositories.ProjectRepository
}

func NewIndexService(injector *do.Injector) (IndexService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &IndexServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
	}, nil
}

func (s *IndexServiceImpl) List(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) ([]string, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return nil, errs.NewForbiddenError("project.error.viewForbidden")
	}

	clientIndexRepo, _, err := s.connectionService.GetClientIndexRepo(project.DBName, nil)
	if err != nil {
		return nil, err
	}

	_, tableName := utils.ParseTableName(fullTableName)
	return clientIndexRepo.List(tableName)
}

func (s *IndexServiceImpl) GetByName(indexName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("project.error.viewForbidden")
	}

	clientIndexRepo, _, err := s.connectionService.GetClientIndexRepo(project.DBName, nil)
	if err != nil {
		return "", err
	}

	_, tableName := utils.ParseTableName(fullTableName)
	return clientIndexRepo.GetByName(tableName, indexName)
}

func (s *IndexServiceImpl) Create(fullTableName string, request *requests.IndexCreateRequest, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return "", errs.NewForbiddenError("table.error.createForbidden")
	}

	clientIndexRepo, _, err := s.connectionService.GetClientIndexRepo(project.DBName, nil)
	if err != nil {
		return "", err
	}

	hasIndex, err := clientIndexRepo.Has(fullTableName, request.Name)
	if err != nil {
		return "", err
	}

	if hasIndex {
		return "", errs.NewUnprocessableError("index.error.alreadyExists")
	}

	_, err = clientIndexRepo.Create(fullTableName, request.Name, request.Columns, request.IsUnique)
	if err != nil {
		return "", err
	}

	_, tableName := utils.ParseTableName(fullTableName)

	return clientIndexRepo.GetByName(tableName, request.Name)
}

func (s *IndexServiceImpl) Delete(indexName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientIndexRepo, _, err := s.connectionService.GetClientIndexRepo(project.DBName, nil)
	if err != nil {
		return false, err
	}

	_, tableName := utils.ParseTableName(fullTableName)
	hasIndex, err := clientIndexRepo.Has(tableName, indexName)
	if err != nil {
		return false, err
	}

	if !hasIndex {
		return false, errs.NewNotFoundError("index.error.notFound")
	}

	return clientIndexRepo.DropIfExists(indexName)
}
