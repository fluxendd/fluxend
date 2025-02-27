package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/column_requests"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ColumnService interface {
	CreateMany(projectID, tableID uuid.UUID, request *column_requests.CreateRequest, authUser models.AuthUser) (models.Table, error)
	AlterMany(tableID, projectID uuid.UUID, request *column_requests.CreateRequest, authUser models.AuthUser) (*models.Table, error)
	Rename(columnName string, tableID, projectID uuid.UUID, request *column_requests.RenameRequest, authUser models.AuthUser) (*models.Table, error)
	Delete(columnName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ColumnServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	projectRepo       *repositories.ProjectRepository
	coreTableRepo     *repositories.CoreTableRepository
}

func NewColumnService(injector *do.Injector) (ColumnService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &ColumnServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
		coreTableRepo:     coreTableRepo,
	}, nil
}

func (s *ColumnServiceImpl) CreateMany(projectID, tableID uuid.UUID, request *column_requests.CreateRequest, authUser models.AuthUser) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return models.Table{}, err
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	anyColumnExists, err := clientColumnRepo.HasAny(table.Name, request.Columns)
	if err != nil {
		return models.Table{}, err
	}

	if anyColumnExists {
		return models.Table{}, errs.NewUnprocessableError("column.error.someAlreadyExist")
	}

	err = clientColumnRepo.CreateMany(table.Name, request.Columns)
	if err != nil {
		return models.Table{}, err
	}

	table.Columns = append(table.Columns, request.Columns...)
	table.UpdatedBy = authUser.Uuid

	_, err = s.coreTableRepo.Update(&table)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *ColumnServiceImpl) AlterMany(tableID, projectID uuid.UUID, request *column_requests.CreateRequest, authUser models.AuthUser) (*models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return &models.Table{}, err
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return &models.Table{}, err
	}

	allColumnsExist, err := clientColumnRepo.HasAll(table.Name, request.Columns)
	if err != nil {
		return &models.Table{}, err
	}

	if !allColumnsExist {
		return &models.Table{}, errs.NewNotFoundError("column.error.someNotFound")
	}

	table.UpdatedBy = authUser.Uuid
	for _, column := range request.Columns {
		for i, tableColumn := range table.Columns {
			if tableColumn.Name == column.Name {
				table.Columns[i].Type = column.Type
				break
			}
		}
	}

	err = clientColumnRepo.AlterMany(table.Name, request.Columns)
	if err != nil {
		return &models.Table{}, err
	}

	return s.coreTableRepo.Update(&table)
}

func (s *ColumnServiceImpl) Rename(columnName string, tableID, projectID uuid.UUID, request *column_requests.RenameRequest, authUser models.AuthUser) (*models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	columnExists, err := s.coreTableRepo.HasColumn(columnName, tableID)
	if err != nil {
		return &models.Table{}, err
	}

	if !columnExists {
		return &models.Table{}, errs.NewNotFoundError("column.error.notFound")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return &models.Table{}, err
	}

	table.UpdatedBy = authUser.Uuid

	for i, column := range table.Columns {
		if column.Name == columnName {
			table.Columns[i].Name = request.Name
			break
		}
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return &models.Table{}, err
	}

	err = clientColumnRepo.Rename(table.Name, columnName, request.Name)
	if err != nil {
		return &models.Table{}, err
	}

	return s.coreTableRepo.Update(&table)
}

func (s *ColumnServiceImpl) Delete(columnName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return false, err
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return false, err
	}

	hasColumn, err := s.coreTableRepo.HasColumn(columnName, tableID)
	if err != nil {
		return false, err
	}

	if !hasColumn {
		return false, errs.NewNotFoundError("column.error.notFound")
	}

	for i, column := range table.Columns {
		if column.Name == columnName {
			// Remove column from slice
			table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
			break
		}
	}

	table.UpdatedBy = authUser.Uuid

	err = clientColumnRepo.Drop(table.Name, columnName)
	if err != nil {
		return false, err
	}

	_, err = s.coreTableRepo.Update(&table)

	return err == nil, err
}

func (s *ColumnServiceImpl) validateNameForDuplication(name string, tableID uuid.UUID) error {
	exists, err := s.coreTableRepo.HasColumn(name, tableID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("table.error.duplicateName")
	}

	return nil
}
