package database

import (
	"errors"
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	flxErrors "fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type TableService interface {
	List(projectUUID uuid.UUID, authUser auth.User) ([]Table, error)
	GetByName(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (Table, error)
	Create(request CreateTableInput, authUser auth.User) (Table, error)
	Upload(request UploadTableInput, authUser auth.User) (Table, error)
	Duplicate(fullTableName string, authUser auth.User, request RenameTableInput) (*Table, error)
	Rename(fullTableName string, authUser auth.User, request RenameTableInput) (Table, error)
	Delete(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type TableServiceImpl struct {
	connectionService ConnectionService
	fileImportService FileImportService
	projectPolicy     *project.Policy
	postgrestService  shared.PostgrestService
	projectRepo       project.Repository
}

func NewTableService(injector *do.Injector) (TableService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)
	fileImportService := do.MustInvoke[FileImportService](injector)
	postgrestService := do.MustInvoke[shared.PostgrestService](injector)

	return &TableServiceImpl{
		connectionService: connectionService,
		fileImportService: fileImportService,
		projectPolicy:     policy,
		projectRepo:       projectRepo,
		postgrestService:  postgrestService,
	}, nil
}

func (s *TableServiceImpl) List(projectUUID uuid.UUID, authUser auth.User) ([]Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return []Table{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return []Table{}, flxErrors.NewForbiddenError("project.error.listForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
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

func (s *TableServiceImpl) GetByName(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("project.error.viewForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return Table{}, err
	}

	return fetchedTable, nil
}

func (s *TableServiceImpl) Create(request CreateTableInput, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("table.error.createForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	if err = s.validateNameForDuplication(request.Name, clientTableRepo); err != nil {
		return Table{}, err
	}

	if err = clientTableRepo.Create(request.Name, request.Columns); err != nil {
		return Table{}, err
	}

	s.postgrestService.RefreshSchemaCache(fetchedProject.DBName)

	return clientTableRepo.GetByNameInSchema(pkg.ParseTableName(request.Name))
}

func (s *TableServiceImpl) Upload(request UploadTableInput, authUser auth.User) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanCreate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("table.error.createForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	if err = s.validateNameForDuplication(request.Name, clientTableRepo); err != nil {
		return Table{}, err
	}

	clientRowRepo, err := s.getClientRowRepo(fetchedProject.DBName, connection)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	file, err := request.File.Open()
	if err != nil {
		return Table{}, err
	}
	defer file.Close()

	columns, values, err := s.fileImportService.ImportCSV(file)
	if err != nil {
		return Table{}, err
	}

	if err = clientTableRepo.Create(request.Name, columns); err != nil {
		return Table{}, err
	}

	if err = clientRowRepo.CreateMany(request.Name, columns, values); err != nil {
		return Table{}, err
	}

	s.postgrestService.RefreshSchemaCache(fetchedProject.DBName)

	return clientTableRepo.GetByNameInSchema(pkg.ParseTableName(request.Name))
}

func (s *TableServiceImpl) Duplicate(fullTableName string, authUser auth.User, request RenameTableInput) (*Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return &Table{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return &Table{}, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return &Table{}, err
	}
	defer connection.Close()

	if err = s.validateNameForDuplication(request.Name, clientTableRepo); err != nil {
		return &Table{}, err
	}

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return &Table{}, err
	}

	if err = clientTableRepo.Duplicate(fetchedTable.Name, request.Name); err != nil {
		return &Table{}, err
	}

	fetchedTable.Name = request.Name
	s.postgrestService.RefreshSchemaCache(fetchedProject.DBName)

	return &fetchedTable, nil
}

func (s *TableServiceImpl) Rename(fullTableName string, authUser auth.User, request RenameTableInput) (Table, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(request.ProjectUUID)
	if err != nil {
		return Table{}, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return Table{}, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return Table{}, err
	}
	defer connection.Close()

	if err = s.validateNameForDuplication(request.Name, clientTableRepo); err != nil {
		return Table{}, err
	}

	fetchedTable, err := clientTableRepo.GetByNameInSchema(pkg.ParseTableName(fullTableName))
	if err != nil {
		return Table{}, err
	}

	if err = clientTableRepo.Rename(fetchedTable.Name, request.Name); err != nil {
		return Table{}, err
	}

	fetchedTable.Name = request.Name

	return fetchedTable, nil
}

func (s *TableServiceImpl) Delete(fullTableName string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return false, flxErrors.NewForbiddenError("project.error.updateForbidden")
	}

	clientTableRepo, connection, err := s.getClientTableRepo(fetchedProject.DBName)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	if err = clientTableRepo.DropIfExists(fullTableName); err != nil {
		return false, err
	}

	return true, nil
}

func (s *TableServiceImpl) getClientTableRepo(dbName string) (TableRepository, *sqlx.DB, error) {
	repo, connection, err := s.connectionService.GetTableRepo(dbName, nil)
	if err != nil {
		return nil, nil, err
	}

	clientRepo, ok := repo.(TableRepository)
	if !ok {
		connection.Close()

		return nil, nil, errors.New("clientTableRepo is not of type *repositories.TableRepository")
	}

	return clientRepo, connection, nil
}

func (s *TableServiceImpl) getClientRowRepo(dbName string, connection *sqlx.DB) (RowRepository, error) {
	repo, _, err := s.connectionService.GetRowRepo(dbName, connection)
	if err != nil {
		return nil, err
	}

	clientRepo, ok := repo.(RowRepository)
	if !ok {
		return nil, errors.New("clientRowRepo is not of type *repositories.RowRepository")
	}

	return clientRepo, nil
}

func (s *TableServiceImpl) validateNameForDuplication(name string, clientTableRepo TableRepository) error {
	exists, err := clientTableRepo.Exists(name)
	if err != nil {
		return err
	}

	if exists {
		return flxErrors.NewUnprocessableError("table.error.alreadyExists")
	}

	return nil
}
