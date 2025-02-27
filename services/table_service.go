package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/table_requests"
	"fluxton/utils"
	"github.com/google/uuid"

	"github.com/samber/do"
)

type TableService interface {
	List(paginationParams utils.PaginationParams, projectID uuid.UUID, authUser models.AuthUser) ([]models.Table, error)
	GetByID(tableID, projectID uuid.UUID, authUser models.AuthUser) (models.Table, error)
	Create(request *table_requests.TableCreateRequest, projectID uuid.UUID, authUser models.AuthUser) (models.Table, error)
	Duplicate(tableID, projectID uuid.UUID, authUser models.AuthUser, request *table_requests.TableRenameRequest) (*models.Table, error)
	Rename(tableID, projectID uuid.UUID, authUser models.AuthUser, request *table_requests.TableRenameRequest) (models.Table, error)
	Delete(tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error)
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

func (s *TableServiceImpl) List(paginationParams utils.PaginationParams, projectID uuid.UUID, authUser models.AuthUser) ([]models.Table, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectID)
	if err != nil {
		return []models.Table{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Table{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.coreTableRepo.ListForProject(paginationParams, projectID)
}

func (s *TableServiceImpl) GetByID(tableID, projectID uuid.UUID, authUser models.AuthUser) (models.Table, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return s.coreTableRepo.GetByID(tableID)
}

func (s *TableServiceImpl) Create(request *table_requests.TableCreateRequest, projectID uuid.UUID, authUser models.AuthUser) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, projectID)
	if err != nil {
		return models.Table{}, err
	}

	// TODO: cleanup name and column names for spaces etc
	table := models.Table{
		Name:        request.Name,
		ProjectUuid: projectID,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
		Columns:     request.Columns,
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

func (s *TableServiceImpl) Duplicate(tableID, projectID uuid.UUID, authUser models.AuthUser, request *table_requests.TableRenameRequest) (*models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.validateNameForDuplication(request.Name, projectID)
	if err != nil {
		return &models.Table{}, err
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return &models.Table{}, err
	}

	clientTableRepo, err := s.getClientTableRepo(project.DBName)
	if err != nil {
		return &models.Table{}, err
	}

	err = clientTableRepo.Duplicate(table.Name, request.Name)
	if err != nil {
		return &models.Table{}, err
	}

	table.Name = request.Name

	return s.coreTableRepo.Create(&table)
}

func (s *TableServiceImpl) Rename(tableID, projectID uuid.UUID, authUser models.AuthUser, request *table_requests.TableRenameRequest) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
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

	return s.coreTableRepo.Rename(tableID, request.Name, authUser.Uuid)
}

func (s *TableServiceImpl) Delete(tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error) {
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
