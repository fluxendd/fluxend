package database

import (
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type IndexService interface {
	List(fullTableName string, projectUUID uuid.UUID, authUser auth.User) ([]string, error)
	GetByName(indexName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (string, error)
	Create(fullTableName string, request CreateIndexInput, authUser auth.User) (string, error)
	Delete(indexName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type IndexServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *project.Policy
	projectRepo       project.Repository
}

func NewIndexService(injector *do.Injector) (IndexService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &IndexServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
	}, nil
}

func (s *IndexServiceImpl) List(fullTableName string, projectUUID uuid.UUID, authUser auth.User) ([]string, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return nil, errors.NewForbiddenError("project.error.viewForbidden")
	}

	clientIndexRepo, connection, err := s.connectionService.GetIndexRepo(fetchedProject.DBName, nil)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	_, tableName := pkg.ParseTableName(fullTableName)
	return clientIndexRepo.List(tableName)
}

func (s *IndexServiceImpl) GetByName(indexName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (string, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return "", errors.NewForbiddenError("project.error.viewForbidden")
	}

	clientIndexRepo, connection, err := s.connectionService.GetIndexRepo(fetchedProject.DBName, nil)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	_, tableName := pkg.ParseTableName(fullTableName)
	return clientIndexRepo.GetByName(tableName, indexName)
}

func (s *IndexServiceImpl) Create(fullTableName string, request CreateIndexInput, authUser auth.User) (string, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return "", errors.NewForbiddenError("table.error.createForbidden")
	}

	clientIndexRepo, connection, err := s.connectionService.GetIndexRepo(fetchedProject.DBName, nil)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	hasIndex, err := clientIndexRepo.Has(fullTableName, request.Name)
	if err != nil {
		return "", err
	}

	if hasIndex {
		return "", errors.NewUnprocessableError("index.error.alreadyExists")
	}

	_, err = clientIndexRepo.Create(fullTableName, request.Name, request.Columns, request.IsUnique)
	if err != nil {
		return "", err
	}

	_, tableName := pkg.ParseTableName(fullTableName)

	return clientIndexRepo.GetByName(tableName, request.Name)
}

func (s *IndexServiceImpl) Delete(indexName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return false, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientIndexRepo, connection, err := s.connectionService.GetIndexRepo(fetchedProject.DBName, nil)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	_, tableName := pkg.ParseTableName(fullTableName)
	hasIndex, err := clientIndexRepo.Has(tableName, indexName)
	if err != nil {
		return false, err
	}

	if !hasIndex {
		return false, errors.NewNotFoundError("index.error.notFound")
	}

	return clientIndexRepo.DropIfExists(indexName)
}
