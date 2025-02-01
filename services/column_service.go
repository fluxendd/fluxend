package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/samber/do"
)

type ColumnService interface {
	Create(request *requests.ColumnCreateRequest, projectID, tableID uint, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	Alter(tableID, projectID uint, authenticatedUser models.AuthenticatedUser, request *requests.TableRenameRequest) (models.Table, error)
	Delete(tableID, organizationID, projectID uint, authenticatedUser models.AuthenticatedUser) (bool, error)
}

type ColumnServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	databaseRepo  *repositories.DatabaseRepository
	projectRepo   *repositories.ProjectRepository
	coreTableRepo *repositories.CoreTableRepository
}

func NewColumnService(injector *do.Injector) (ColumnService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &ColumnServiceImpl{
		projectPolicy: policy,
		databaseRepo:  databaseRepo,
		projectRepo:   projectRepo,
		coreTableRepo: coreTableRepo,
	}, nil
}

func (s *ColumnServiceImpl) Create(request *requests.ColumnCreateRequest, projectID, tableID uint, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
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

	_, err = s.coreTableRepo.Update(&table)
	if err != nil {
		return models.Table{}, err
	}

	clientColumnRepo, err := s.getClientColumnRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	_, err = clientColumnRepo.Create(table.Name, request.Column)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *ColumnServiceImpl) Alter(tableID, projectID uint, authenticatedUser models.AuthenticatedUser, request *requests.TableRenameRequest) (models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(request.OrganizationID, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.validateNameForDuplication(request.Name, projectID)
	if err != nil {
		return models.Table{}, err
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return models.Table{}, err
	}

	clientTableRepo, err := s.getClientTableRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	err = clientTableRepo.Rename(table.Name, request.Name)
	if err != nil {
		return models.Table{}, err
	}

	return s.coreTableRepo.Rename(tableID, request.Name)
}

func (s *ColumnServiceImpl) Delete(tableID, organizationID, projectID uint, authenticatedUser models.AuthenticatedUser) (bool, error) {
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

	clientTableRepo, err := s.getClientTableRepo(project.DBName)
	if err != nil {
		return false, err
	}

	err = clientTableRepo.DropIfExists(table.Name)
	if err != nil {
		return false, err
	}

	return s.coreTableRepo.Delete(tableID)
}

func (s *ColumnServiceImpl) validateNameForDuplication(name string, tableID uint) error {
	exists, err := s.coreTableRepo.HasColumn(name, tableID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("table.error.duplicateName")
	}

	return nil
}

func (s *ColumnServiceImpl) getClientTableRepo(databaseName string) (*repositories.ClientTableRepository, error) {
	clientDatabaseConnection, err := s.databaseRepo.Connect(databaseName)
	if err != nil {
		return nil, err
	}

	clientTableRepo, err := repositories.NewClientTableRepository(clientDatabaseConnection)
	if err != nil {
		return nil, err
	}

	return clientTableRepo, nil
}

func (s *ColumnServiceImpl) getClientColumnRepo(databaseName string) (*repositories.ClientColumnRepository, error) {
	clientDatabaseConnection, err := s.databaseRepo.Connect(databaseName)
	if err != nil {
		return nil, err
	}

	clientColumnRepo, err := repositories.NewClientColumnRepository(clientDatabaseConnection)
	if err != nil {
		return nil, err
	}

	return clientColumnRepo, nil
}
