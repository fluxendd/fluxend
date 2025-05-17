package table

import (
	"fluxton/internal/adapters/connection"
	table2 "fluxton/internal/api/dto/database/table"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/import"
	"fluxton/internal/domain/project"
	"fluxton/models"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type TableService interface {
	List(projectUUID uuid.UUID, authUser models.AuthUser) ([]Table, error)
	GetByName(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (Table, error)
	Create(request *table2.CreateRequest, authUser models.AuthUser) (Table, error)
	Upload(request *table2.UploadRequest, authUser models.AuthUser) (Table, error)
	Duplicate(fullTableName string, authUser models.AuthUser, request *table2.RenameRequest) (*Table, error)
	Rename(fullTableName string, authUser models.AuthUser, request *table2.RenameRequest) (Table, error)
	Delete(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type TableServiceImpl struct {
	connectionService connection.ConnectionService
	fileImportService _import.FileImportService
	projectPolicy     *project.ProjectPolicy
	databaseRepo      *repositories2.DatabaseRepository
	projectRepo       *repositories2.ProjectRepository
}

func NewTableService(injector *do.Injector) (TableService, error) {
	connectionService := do.MustInvoke[connection.ConnectionService](injector)
	policy := do.MustInvoke[*project.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories2.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)
	fileImportService := do.MustInvoke[_import.FileImportService](injector)

	return &TableServiceImpl{
		connectionService: connectionService,
		fileImportService: fileImportService,
		projectPolicy:     policy,
		databaseRepo:      databaseRepo,
		projectRepo:       projectRepo,
	}, nil
}

func (s *TableServiceImpl) List(projectUUID uuid.UUID, authUser models.AuthUser) ([]Table, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return []Table{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return []Table{}, errors.NewForbiddenError("project.error.listForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return []Table{}, err
	}
	defer connection.Close()

	tables, err := clientTableRepo.List()
	if err != nil {
		return []Table{}, err
	}

	return tables, nil
}

func (s *TableServiceImpl) GetByName(fullTableName string, projectUUID uuid.UUID, authUser models.AuthUser) (Table, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return Table{}, errors.NewForbiddenError("project.error.viewForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return table.Table{}, err
	}

	return table, nil
}

func (s *TableServiceImpl) Create(request *table2.CreateRequest, authUser models.AuthUser) (Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return Table{}, errors.NewForbiddenError("table.error.createForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	err = s.validateNameForDuplication(request.Name, clientTableRepo)
	if err != nil {
		return Table{}, err
	}

	err = clientTableRepo.Create(request.Name, request.Columns)
	if err != nil {
		return Table{}, err
	}

	return clientTableRepo.GetByNameInSchema(pkg.ParseTableName(request.Name))
}

func (s *TableServiceImpl) Upload(request *table2.UploadRequest, authUser models.AuthUser) (Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(project.OrganizationUuid, authUser) {
		return Table{}, errors.NewForbiddenError("table.error.createForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	err = s.validateNameForDuplication(request.Name, clientTableRepo)
	if err != nil {
		return Table{}, err
	}

	clientRowRepo, _, err := s.connectionService.GetRowRepo(project.DBName, connection)

	file, err := request.File.Open()
	if err != nil {
		return Table{}, err
	}
	defer file.Close()

	columns, values, err := s.fileImportService.ImportCSV(file)
	if err != nil {
		return Table{}, err
	}

	err = clientTableRepo.Create(request.Name, columns)
	if err != nil {
		return Table{}, err
	}

	err = clientRowRepo.CreateMany(request.Name, columns, values)
	if err != nil {
		return Table{}, err
	}

	return clientTableRepo.GetByNameInSchema(pkg.ParseTableName(request.Name))
}

func (s *TableServiceImpl) Duplicate(fullTableName string, authUser models.AuthUser, request *table2.RenameRequest) (*Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return &Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &Table{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return &Table{}, err
	}
	defer connection.Close()

	err = s.validateNameForDuplication(request.Name, clientTableRepo)
	if err != nil {
		return &Table{}, err
	}

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return &table.Table{}, err
	}

	err = clientTableRepo.Duplicate(table.Name, request.Name)
	if err != nil {
		return &table.Table{}, err
	}

	table.Name = request.Name

	return &table, nil
}

func (s *TableServiceImpl) Rename(fullTableName string, authUser models.AuthUser, request *table2.RenameRequest) (Table, error) {
	project, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return Table{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	err = s.validateNameForDuplication(request.Name, clientTableRepo)
	if err != nil {
		return Table{}, err
	}

	table, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return table.Table{}, err
	}

	err = clientTableRepo.Rename(table.Name, request.Name)
	if err != nil {
		return table.Table{}, err
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
		return false, errors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.connectionService.GetTableRepo(project.DBName, nil)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	err = clientTableRepo.DropIfExists(fullTableName)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *TableServiceImpl) validateNameForDuplication(name string, clientTableRepo *repositories2.TableRepository) error {
	exists, err := clientTableRepo.Exists(name)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("table.error.alreadyExists")
	}

	return nil
}
