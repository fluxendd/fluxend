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
	List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Table, error)
	GetByID(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (models.Table, error)
	Create(request *table_requests.CreateRequest, authUser models.AuthUser) (models.Table, error)
	Duplicate(fullTableName string, authUser models.AuthUser, request *table_requests.RenameRequest) (*models.Table, error)
	Rename(fullTableName string, authUser models.AuthUser, request *table_requests.RenameRequest) (models.Table, error)
	Delete(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type TableServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	databaseRepo      *repositories.DatabaseRepository
	projectRepo       *repositories.ProjectRepository
}

func NewTableService(injector *do.Injector) (TableService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &TableServiceImpl{
		connectionService: connectionService,
		projectPolicy:     policy,
		databaseRepo:      databaseRepo,
		projectRepo:       projectRepo,
	}, nil
}

func (s *TableServiceImpl) List(projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return []models.Table{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return []models.Table{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return []models.Table{}, err
	}

	tables, err := clientTableRepo.List()
	if err != nil {
		return []models.Table{}, err
	}

	return tables, nil
}

func (s *TableServiceImpl) GetByID(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return models.Table{}, err
	}

	tableName, schema := utils.ParseTableName(fullTableName)
	table, err := clientTableRepo.GetByNameInSchema(tableName, schema)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *TableServiceImpl) Create(request *table_requests.CreateRequest, authUser models.AuthUser) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	/*err = s.validateNameForDuplication(request.Name, request.ProjectUUID)
	if err != nil {
		return models.Table{}, err
	}*/

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return models.Table{}, err
	}

	err = clientTableRepo.Create(request.Name, request.Columns)
	if err != nil {
		return models.Table{}, err
	}

	return models.Table{}, nil
}

func (s *TableServiceImpl) Duplicate(fullTableName string, authUser models.AuthUser, request *table_requests.RenameRequest) (*models.Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return &models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	/*err = s.validateNameForDuplication(request.Name, project.Uuid)
	if err != nil {
		return &models.Table{}, err
	}*/

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return &models.Table{}, err
	}

	tableName, schema := utils.ParseTableName(fullTableName)
	table, err := clientTableRepo.GetByNameInSchema(tableName, schema)
	if err != nil {
		return &models.Table{}, err
	}

	err = clientTableRepo.Duplicate(table.Name, request.Name)
	if err != nil {
		return &models.Table{}, err
	}

	table.Name = request.Name

	return &table, nil
}

func (s *TableServiceImpl) Rename(fullTableName string, authUser models.AuthUser, request *table_requests.RenameRequest) (models.Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return models.Table{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	/*err = s.validateNameForDuplication(request.Name, project.Uuid)
	if err != nil {
		return models.Table{}, err
	}*/

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return models.Table{}, err
	}

	tableName, schema := utils.ParseTableName(fullTableName)
	table, err := clientTableRepo.GetByNameInSchema(tableName, schema)
	if err != nil {
		return models.Table{}, err
	}

	err = clientTableRepo.Rename(table.Name, request.Name)
	if err != nil {
		return models.Table{}, err
	}

	table.Name = request.Name

	return table, nil
}

func (s *TableServiceImpl) Delete(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, _, err := s.connectionService.GetClientTableRepo(project.DBName, nil)
	if err != nil {
		return false, err
	}

	err = clientTableRepo.DropIfExists(fullTableName)
	if err != nil {
		return false, err
	}

	return true, nil
}

/*func (s *TableServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.coreTableRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("table.error.duplicateName")
	}

	return nil
}*/
