package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	//"fluxton/utils"
	//"github.com/google/uuid"
	"github.com/samber/do"
	//"strings"
)

type TableService interface {
	//List(paginationParams utils.PaginationParams, organizationId, authenticatedUserId uint) ([]models.Table, error)
	//GetByID(tableId, organizationId uint, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	Create(request *requests.TableCreateRequest, projectID uint, authenticatedUser models.AuthenticatedUser) (models.Table, error)
	//Update(tableId uint, authenticatedUser models.AuthenticatedUser, request *requests.TableCreateRequest) (*models.Table, error)
	Delete(tableId uint, authenticatedUser models.AuthenticatedUser) (bool, error)
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

/*func (s *TableServiceImpl) List(paginationParams utils.PaginationParams, organizationId, authenticatedUserId uint) ([]models.Table, error) {
	if !s.projectPolicy.CanList(organizationId, authenticatedUserId) {
		return []models.Table{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authenticatedUserId)
}

func (s *TableServiceImpl) GetByID(tableId, organizationId uint, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
	if !s.projectPolicy.CanView(organizationId, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return s.projectRepo.GetByID(tableId)
}*/

func (s *TableServiceImpl) Create(request *requests.TableCreateRequest, projectID uint, authenticatedUser models.AuthenticatedUser) (models.Table, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationID, authenticatedUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	err := s.validateNameForDuplication(request.Name, request.OrganizationID)
	if err != nil {
		return models.Table{}, err
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	table := models.Table{
		Name:      request.Name,
		ProjectID: projectID,
		Fields:    request.Fields,
	}

	_, err = s.coreTableRepo.Create(&table)
	if err != nil {
		return models.Table{}, err
	}

	clientDatabaseConnection, err := s.databaseRepo.Connect(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	clientTableRepo, err := repositories.NewClientTableRepository(clientDatabaseConnection)
	if err != nil {
		return models.Table{}, err
	}

	err = clientTableRepo.Create(table.Name, table.Fields)

	return table, nil
}

/*func (s *TableServiceImpl) Update(tableId uint, authenticatedUser models.AuthenticatedUser, request *requests.TableCreateRequest) (*models.Table, error) {
	project, err := s.projectRepo.GetByID(tableId)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(tableId, authenticatedUser) {
		return &models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	project.OrganizationID = request.OrganizationID
	err = utils.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	err = s.validateNameForDuplication(request.Name, request.OrganizationID)
	if err != nil {
		return &models.Table{}, err
	}

	return s.projectRepo.Update(tableId, &project)
}*/

func (s *TableServiceImpl) Delete(tableId uint, authenticatedUser models.AuthenticatedUser) (bool, error) {
	if !s.projectPolicy.CanUpdate(tableId, authenticatedUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	return s.projectRepo.Delete(tableId)
}

func (s *TableServiceImpl) validateNameForDuplication(name string, organizationId uint) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationId)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
