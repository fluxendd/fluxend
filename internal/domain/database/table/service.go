package table

import (
	"errors"
	"fluxton/internal/api/dto/database/table"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/database/client"
	"fluxton/internal/domain/database/row"
	"fluxton/internal/domain/file_import"
	"fluxton/internal/domain/project"
	"fluxton/pkg"
	flxErrors "fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Service interface {
	List(projectUUID uuid.UUID, authUser auth.User) ([]Table, error)
	GetByName(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (Table, error)
	Create(request *table.CreateRequest, authUser auth.User) (Table, error)
	Upload(request *table.UploadRequest, authUser auth.User) (Table, error)
	Duplicate(fullTableName string, authUser auth.User, request *table.RenameRequest) (*Table, error)
	Rename(fullTableName string, authUser auth.User, request *table.RenameRequest) (Table, error)
	Delete(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	connectionService client.ConnectionService
	fileImportService file_import.Service
	projectPolicy     *project.Policy
	databaseRepo      client.DatabaseService
	projectRepo       project.Repository
}

func NewTableService(injector *do.Injector) (Service, error) {
	connectionService := do.MustInvoke[client.ConnectionService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	databaseRepo := do.MustInvoke[client.DatabaseService](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)
	fileImportService := do.MustInvoke[file_import.Service](injector)

	return &ServiceImpl{
		connectionService: connectionService,
		fileImportService: fileImportService,
		projectPolicy:     policy,
		databaseRepo:      databaseRepo,
		projectRepo:       projectRepo,
	}, nil
}

func (s *ServiceImpl) List(projectUUID uuid.UUID, authUser auth.User) ([]Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return []Table{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return []Table{}, flxErrors.NewForbiddenError("project.error.listForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return []Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return []Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	tables, err := clientTableRepo.List()
	if err != nil {
		return []Table{}, err
	}

	return tables, nil
}

func (s *ServiceImpl) GetByName(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("project.error.viewForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return Table{}, err
	}

	return fetchedTable, nil
}

func (s *ServiceImpl) Create(request *table.CreateRequest, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("table.error.createForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	err = s.validateNameForDuplication(request.Name)
	if err != nil {
		return Table{}, err
	}

	err = clientTableRepo.Create(request.Name, request.Columns)
	if err != nil {
		return Table{}, err
	}

	return clientTableRepo.GetByNameInSchema(pkg.ParseTableName(request.Name))
}

func (s *ServiceImpl) Upload(request *table.UploadRequest, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("table.error.createForbidden")
	}

	fetchedTableRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedTableRepo.(Repository)
	if !ok {
		return Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	err = s.validateNameForDuplication(request.Name)
	if err != nil {
		return Table{}, err
	}

	fetchedRowRepo, _, err := s.connectionService.GetRowRepo(fetchedProject.DBName, connection)
	if err != nil {
		return Table{}, err
	}

	clientRowRepo, ok := fetchedRowRepo.(row.Repository)
	if !ok {
		return Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

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

func (s *ServiceImpl) Duplicate(fullTableName string, authUser auth.User, request *table.RenameRequest) (*Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return &Table{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return &Table{}, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return &Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return &Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	err = s.validateNameForDuplication(request.Name)
	if err != nil {
		return &Table{}, err
	}

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return &Table{}, err
	}

	err = clientTableRepo.Duplicate(fetchedTable.Name, request.Name)
	if err != nil {
		return &Table{}, err
	}

	fetchedTable.Name = request.Name

	return &fetchedTable, nil
}

func (s *ServiceImpl) Rename(fullTableName string, authUser auth.User, request *table.RenameRequest) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return Table{}, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	err = s.validateNameForDuplication(request.Name)
	if err != nil {
		return Table{}, err
	}

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return Table{}, err
	}

	err = clientTableRepo.Rename(fetchedTable.Name, request.Name)
	if err != nil {
		return Table{}, err
	}

	fetchedTable.Name = request.Name

	return fetchedTable, nil
}

func (s *ServiceImpl) Delete(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return false, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	fetchedRepo, connection, err := s.connectionService.GetTableRepo(fetchedProject.DBName, nil)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	clientTableRepo, ok := fetchedRepo.(Repository)
	if !ok {
		return false, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	err = clientTableRepo.DropIfExists(fullTableName)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ServiceImpl) validateNameForDuplication(name string) error {
	// TODO: COME_BACK_FOR_ME
	/*exists, err := clientTableRepo.Exists(name)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("table.error.alreadyExists")
	}*/

	return nil
}
