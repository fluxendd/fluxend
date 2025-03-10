package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/column_requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ColumnService interface {
	List(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Column, error)
	CreateMany(fullTableName string, request *column_requests.CreateRequest, authUser models.AuthUser) ([]models.Column, error)
	AlterMany(fullTableName string, request *column_requests.CreateRequest, authUser models.AuthUser) ([]models.Column, error)
	Rename(columnName, fullTableName string, request *column_requests.RenameRequest, authUser models.AuthUser) ([]models.Column, error)
	Delete(columnName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ColumnServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	projectRepo       *repositories.ProjectRepository
}

func NewColumnService(injector *do.Injector) (ColumnService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &ColumnServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
	}, nil
}

func (s *ColumnServiceImpl) List(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Column, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return nil, errs.NewForbiddenError("project.error.readForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return nil, err
	}

	table, err := clientTableRepo.GetByNameInSchema(utils.ParseTableName(fullTableName))
	if err != nil {
		return nil, err
	}

	clientColumnRepo, _, err := s.connectionService.GetClientColumnRepo(project.DBName, connection)
	if err != nil {
		return nil, err
	}

	columns, err := clientColumnRepo.List(table.Name)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

func (s *ColumnServiceImpl) CreateMany(fullTableName string, request *column_requests.CreateRequest, authUser models.AuthUser) ([]models.Column, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []models.Column{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return []models.Column{}, errs.NewForbiddenError("column.error.createForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return []models.Column{}, err
	}

	table, err := clientTableRepo.GetByNameInSchema(utils.ParseTableName(fullTableName))
	if err != nil {
		return []models.Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetClientColumnRepo(project.DBName, connection)
	if err != nil {
		return []models.Column{}, err
	}

	anyColumnExists, err := clientColumnRepo.HasAny(table.Name, request.Columns)
	if err != nil {
		return []models.Column{}, err
	}

	if anyColumnExists {
		return []models.Column{}, errs.NewUnprocessableError("column.error.someAlreadyExist")
	}

	err = clientColumnRepo.CreateMany(table.Name, request.Columns)
	if err != nil {
		return []models.Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) AlterMany(fullTableName string, request *column_requests.CreateRequest, authUser models.AuthUser) ([]models.Column, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []models.Column{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return []models.Column{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return []models.Column{}, err
	}

	table, err := clientTableRepo.GetByNameInSchema(utils.ParseTableName(fullTableName))
	if err != nil {
		return []models.Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetClientColumnRepo(project.DBName, connection)
	if err != nil {
		return []models.Column{}, err
	}

	allColumnsExist, err := clientColumnRepo.HasAll(table.Name, request.Columns)
	if err != nil {
		return []models.Column{}, err
	}

	if !allColumnsExist {
		return []models.Column{}, errs.NewNotFoundError("column.error.someNotFound")
	}

	err = clientColumnRepo.AlterMany(table.Name, request.Columns)
	if err != nil {
		return []models.Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) Rename(columnName string, fullTableName string, request *column_requests.RenameRequest, authUser models.AuthUser) ([]models.Column, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return []models.Column{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return []models.Column{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return []models.Column{}, err
	}

	table, err := clientTableRepo.GetByNameInSchema(utils.ParseTableName(fullTableName))
	if err != nil {
		return []models.Column{}, err
	}

	clientColumnRepo, _, err := s.connectionService.GetClientColumnRepo(project.DBName, connection)
	if err != nil {
		return []models.Column{}, err
	}

	columnExists, err := clientColumnRepo.Has(table.Name, columnName)
	if err != nil {
		return []models.Column{}, err
	}

	if !columnExists {
		return []models.Column{}, errs.NewNotFoundError("column.error.notFound")
	}

	err = clientColumnRepo.Rename(table.Name, columnName, request.Name)
	if err != nil {
		return []models.Column{}, err
	}

	return clientColumnRepo.List(table.Name)
}

func (s *ColumnServiceImpl) Delete(columnName, fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientColumnRepo, _, err := s.connectionService.GetClientColumnRepo(project.DBName, nil)
	if err != nil {
		return false, err
	}

	_, tableName := utils.ParseTableName(fullTableName)
	columnExists, err := clientColumnRepo.Has(tableName, columnName)
	if err != nil {
		return false, err
	}

	if !columnExists {
		return false, errs.NewNotFoundError("column.error.notFound")
	}

	err = clientColumnRepo.Drop(fullTableName, columnName)
	if err != nil {
		return false, err
	}

	return true, err
}
