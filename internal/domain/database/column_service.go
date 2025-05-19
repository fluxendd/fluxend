package database

import (
	"fluxton/internal/api/dto/database/column"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ColumnService interface {
	List(fullTableName string, projectUUID uuid.UUID, authUser auth.User) ([]Column, error)
	CreateMany(fullTableName string, request *column.CreateRequest, authUser auth.User) ([]Column, error)
	AlterMany(fullTableName string, request *column.CreateRequest, authUser auth.User) ([]Column, error)
	Rename(columnName, fullTableName string, request *column.RenameRequest, authUser auth.User) ([]Column, error)
	Delete(columnName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ColumnServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *project.Policy
	projectRepo       project.Repository
}

func NewColumnService(injector *do.Injector) (ColumnService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ColumnServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
	}, nil
}

func (s *ColumnServiceImpl) List(fullTableName string, projectUUID uuid.UUID, authUser auth.User) ([]Column, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return nil, errors.NewForbiddenError("project.error.viewForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return nil, err
	}

	clientColumnRepo, _, err := s.connectionService.GetColumnRepo(fetchedProject.DBName, connection)
	if err != nil {
		return nil, err
	}

	columns, err := clientColumnRepo.List(table.Name)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

func (s *ColumnServiceImpl) CreateMany(fullTableName string, request *column.CreateRequest, authUser auth.User) ([]Column, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []Column{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return []Column{}, errors.NewForbiddenError("column.error.createForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return []Column{}, err
	}
	defer connection.Close()

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return []Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetColumnRepo(fetchedProject.DBName, connection)
	if err != nil {
		return []Column{}, err
	}

	anyColumnExists, err := clientColumnRepo.HasAny(table.Name, request.Columns)
	if err != nil {
		return []Column{}, err
	}

	if anyColumnExists {
		return []Column{}, errors.NewUnprocessableError("column.error.someAlreadyExist")
	}

	err = clientColumnRepo.CreateMany(table.Name, request.Columns)
	if err != nil {
		return []Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) AlterMany(fullTableName string, request *column.CreateRequest, authUser auth.User) ([]Column, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []Column{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return []Column{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return []Column{}, err
	}
	defer connection.Close()

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return []Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetColumnRepo(fetchedProject.DBName, connection)
	if err != nil {
		return []Column{}, err
	}

	allColumnsExist, err := clientColumnRepo.HasAll(table.Name, request.Columns)
	if err != nil {
		return []Column{}, err
	}

	if !allColumnsExist {
		return []Column{}, errors.NewNotFoundError("column.error.someNotFound")
	}

	err = clientColumnRepo.AlterMany(table.Name, request.Columns)
	if err != nil {
		return []Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) Rename(columnName string, fullTableName string, request *column.RenameRequest, authUser auth.User) ([]Column, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []Column{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return []Column{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return []Column{}, err
	}
	defer connection.Close()

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return []Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetColumnRepo(fetchedProject.DBName, connection)
	if err != nil {
		return []Column{}, err
	}

	columnExists, err := clientColumnRepo.Has(table.Name, columnName)
	if err != nil {
		return []Column{}, err
	}

	if !columnExists {
		return []Column{}, errors.NewNotFoundError("column.error.notFound")
	}

	err = clientColumnRepo.Rename(table.Name, columnName, request.Name)
	if err != nil {
		return []Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) Delete(columnName, fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return false, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientColumnRepo, connection, err := s.connectionService.GetColumnRepo(fetchedProject.DBName, nil)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	_, tableName := pkg.ParseTableName(fullTableName)
	columnExists, err := clientColumnRepo.Has(tableName, columnName)
	if err != nil {
		return false, err
	}

	if !columnExists {
		return false, errors.NewNotFoundError("column.error.notFound")
	}

	err = clientColumnRepo.Drop(fullTableName, columnName)
	if err != nil {
		return false, err
	}

	return true, err
}
