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

type ColumnService interface {
	Create(projectID, tableID uuid.UUID, request *requests.ColumnCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	Alter(columnName string, tableID, projectID uuid.UUID, request *requests.ColumnAlterRequest, authenticatedUser models.AuthenticatedUser) (*models.Table, error)
	Rename(columnName string, tableID, projectID uuid.UUID, request *requests.ColumnRenameRequest, authenticatedUser models.AuthenticatedUser) (*models.Table, error)
	Delete(columnName string, tableID, organizationID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error)
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

func (s *ColumnServiceImpl) Create(projectID, tableID uuid.UUID, request *requests.ColumnCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(request.OrganizationID, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return models.Table{}, err
	}

	err = s.validateNameForDuplication(request.Column.Name, tableID)
	if err != nil {
		return models.Table{}, err
	}

	table.Columns = append(table.Columns, request.Column)
	table.UpdatedBy = authenticatedUser.ID

	_, err = s.coreTableRepo.Update(&table)
	if err != nil {
		return models.Table{}, err
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	err = clientColumnRepo.Create(table.Name, request.Column)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *ColumnServiceImpl) Alter(columnName string, tableID, projectID uuid.UUID, request *requests.ColumnAlterRequest, authenticatedUser models.AuthenticatedUser) (*models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(request.OrganizationID, authenticatedUser) {
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

	table.UpdatedBy = authenticatedUser.ID

	for i, column := range table.Columns {
		if column.Name == columnName {
			table.Columns[i].Type = request.Type
			break
		}
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return &models.Table{}, err
	}

	err = clientColumnRepo.Alter(table.Name, columnName, request.Type)
	if err != nil {
		return &models.Table{}, err
	}

	return s.coreTableRepo.Update(&table)
}

func (s *ColumnServiceImpl) Rename(columnName string, tableID, projectID uuid.UUID, request *requests.ColumnRenameRequest, authenticatedUser models.AuthenticatedUser) (*models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(request.OrganizationID, authenticatedUser) {
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

	table.UpdatedBy = authenticatedUser.ID

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

func (s *ColumnServiceImpl) Delete(columnName string, tableID, organizationID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationID, authenticatedUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return false, err
	}

	for i, column := range table.Columns {
		if column.Name == columnName {
			// Remove column from slice
			table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
			break
		}
	}

	table.UpdatedBy = authenticatedUser.ID

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return false, err
	}

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
