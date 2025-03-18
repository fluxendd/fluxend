package services

import (
	"fluxton/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type ConnectionService interface {
	ConnectByDatabaseName(name string) (*sqlx.DB, error)
	ConnectByProjectUUID(projectUUID uuid.UUID) (*sqlx.DB, error)
	GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (*repositories.DatabaseStatsRepository, *sqlx.DB, error)
	GetTableRepo(databaseName string, connection *sqlx.DB) (*repositories.TableRepository, *sqlx.DB, error)
	GetColumnRepo(databaseName string, connection *sqlx.DB) (*repositories.ColumnRepository, *sqlx.DB, error)
	GetIndexRepo(databaseName string, connection *sqlx.DB) (*repositories.IndexRepository, *sqlx.DB, error)
	GetFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories.FunctionRepository, *sqlx.DB, error)
}

type ConnectionServiceImpl struct {
	databaseRepo *repositories.DatabaseRepository
	projectRepo  *repositories.ProjectRepository
}

func NewConnectionService(injector *do.Injector) (ConnectionService, error) {
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &ConnectionServiceImpl{
		databaseRepo: databaseRepo,
		projectRepo:  projectRepo,
	}, nil
}

func (s *ConnectionServiceImpl) ConnectByDatabaseName(name string) (*sqlx.DB, error) {
	return s.databaseRepo.Connect(name)
}

func (s *ConnectionServiceImpl) ConnectByProjectUUID(projectUUID uuid.UUID) (*sqlx.DB, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	return s.databaseRepo.Connect(project.DBName)
}

func (s *ConnectionServiceImpl) GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (*repositories.DatabaseStatsRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientDatabaseStatsRepo, err := repositories.NewDatabaseStatsRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientDatabaseStatsRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetTableRepo(databaseName string, connection *sqlx.DB) (*repositories.TableRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientTableRepo, err := repositories.NewTableRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientTableRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetColumnRepo(databaseName string, connection *sqlx.DB) (*repositories.ColumnRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientColumnRepo, err := repositories.NewColumnRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientColumnRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetIndexRepo(databaseName string, connection *sqlx.DB) (*repositories.IndexRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientIndexRepo, err := repositories.NewIndexRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientIndexRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories.FunctionRepository, *sqlx.DB, error) {
	databaseName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return nil, nil, err
	}

	return s.getClientFunctionRepo(databaseName, connection)
}

func (s *ConnectionServiceImpl) getClientFunctionRepo(databaseName string, connection *sqlx.DB) (*repositories.FunctionRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientFunctionRepo, err := repositories.NewFunctionRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientFunctionRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) getOrCreateConnection(databaseName string, connection *sqlx.DB) (*sqlx.DB, error) {
	if connection != nil {
		return connection, nil
	}

	return s.databaseRepo.Connect(databaseName)
}
