package client

import (
	repositories2 "fluxton/internal/database/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type Service interface {
	ConnectByDatabaseName(name string) (*sqlx.DB, error)
	ConnectByProjectUUID(projectUUID uuid.UUID) (*sqlx.DB, error)
	GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (*repositories2.DatabaseStatsRepository, *sqlx.DB, error)
	GetTableRepo(databaseName string, connection *sqlx.DB) (*repositories2.TableRepository, *sqlx.DB, error)
	GetColumnRepo(databaseName string, connection *sqlx.DB) (*repositories2.ColumnRepository, *sqlx.DB, error)
	GetIndexRepo(databaseName string, connection *sqlx.DB) (*repositories2.IndexRepository, *sqlx.DB, error)
	GetRowRepo(databaseName string, connection *sqlx.DB) (*repositories2.RowRepository, *sqlx.DB, error)
	GetFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories2.FunctionRepository, *sqlx.DB, error)
}

type ServiceImpl struct {
	databaseRepo *Repository
	projectRepo  *repositories2.ProjectRepository
}

func NewClientService(injector *do.Injector) (Service, error) {
	databaseRepo := do.MustInvoke[*Repository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &ServiceImpl{
		databaseRepo: databaseRepo,
		projectRepo:  projectRepo,
	}, nil
}

func (s *ServiceImpl) ConnectByDatabaseName(name string) (*sqlx.DB, error) {
	return s.databaseRepo.Connect(name)
}

func (s *ServiceImpl) ConnectByProjectUUID(projectUUID uuid.UUID) (*sqlx.DB, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	return s.databaseRepo.Connect(project.DBName)
}

func (s *ServiceImpl) GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (*repositories2.DatabaseStatsRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientDatabaseStatsRepo, err := repositories2.NewDatabaseStatsRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientDatabaseStatsRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetTableRepo(databaseName string, connection *sqlx.DB) (*repositories2.TableRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientTableRepo, err := repositories2.NewTableRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientTableRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetColumnRepo(databaseName string, connection *sqlx.DB) (*repositories2.ColumnRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientColumnRepo, err := repositories2.NewColumnRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientColumnRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetIndexRepo(databaseName string, connection *sqlx.DB) (*repositories2.IndexRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientIndexRepo, err := repositories2.NewIndexRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientIndexRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetRowRepo(databaseName string, connection *sqlx.DB) (*repositories2.RowRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientIndexRepo, err := repositories2.NewRowRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientIndexRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories2.FunctionRepository, *sqlx.DB, error) {
	databaseName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return nil, nil, err
	}

	return s.getClientFunctionRepo(databaseName, connection)
}

func (s *ServiceImpl) getClientFunctionRepo(databaseName string, connection *sqlx.DB) (*repositories2.FunctionRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientFunctionRepo, err := repositories2.NewFunctionRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientFunctionRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) getOrCreateConnection(databaseName string, connection *sqlx.DB) (*sqlx.DB, error) {
	if connection != nil {
		return connection, nil
	}

	return s.databaseRepo.Connect(databaseName)
}
