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

type TableService interface {
	List(paginationParams utils.PaginationParams, organizationID, projectID, authenticatedUserID uuid.UUID) ([]models.Table, error)
	GetByID(tableID, organizationID uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	Create(request *requests.TableCreateRequest, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	Rename(tableID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.TableRenameRequest) (models.Table, error)
	Delete(tableID, organizationID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error)
}

type TableServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	databaseRepo  *repositories.DatabaseRepository
	projectRepo   *repositories.ProjectRepository
	coreTableRepo *repositories.CoreTableRepository
}

func NewTableService(injector *do.Injector) (TableService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &TableServiceImpl{
		projectPolicy: policy,
		databaseRepo:  databaseRepo,
		projectRepo:   projectRepo,
		coreTableRepo: coreTableRepo,
	}, nil
}

func (s *TableServiceImpl) List(paginationParams utils.PaginationParams, organizationID, projectID, authenticatedUserID uuid.UUID) ([]models.Table, error) {
	if !s.projectPolicy.CanList(organizationID, authenticatedUserID) {
		return []models.Table{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.coreTableRepo.ListForProject(paginationParams, projectID)
}

func (s *TableServiceImpl) GetByID(tableID, organizationID uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
	if !s.projectPolicy.CanView(organizationID, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return s.coreTableRepo.GetByID(tableID)
}

func (s *TableServiceImpl) Create(request *requests.TableCreateRequest, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(request.OrganizationID, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, projectID)
	if err != nil {
		return models.Table{}, err
	}

	// TODO: cleanup name and column names for spaces etc
	table := models.Table{
		Name:      request.Name,
		ProjectID: projectID,
		CreatedBy: authenticatedUser.ID,
		UpdatedBy: authenticatedUser.ID,
		Columns:   request.Columns,
	}

	_, err = s.coreTableRepo.Create(&table)
	if err != nil {
		return models.Table{}, err
	}

	clientTableRepo, err := s.getClientTableRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	err = clientTableRepo.Create(table.Name, table.Columns)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *TableServiceImpl) Rename(tableID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.TableRenameRequest) (models.Table, error) {
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

	return s.coreTableRepo.Rename(tableID, request.Name, authenticatedUser.ID)
}

func (s *TableServiceImpl) Delete(tableID, organizationID, projectID uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error) {
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

func (s *TableServiceImpl) validateNameForDuplication(name string, projectID uuid.UUID) error {
	exists, err := s.coreTableRepo.ExistsByNameForProject(name, projectID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("table.error.duplicateName")
	}

	return nil
}

func (s *TableServiceImpl) getClientTableRepo(databaseName string) (*repositories.ClientTableRepository, error) {
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
