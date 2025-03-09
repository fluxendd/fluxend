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
	GetClientTableRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientTableRepository, *sqlx.DB, error)
	GetClientColumnRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientColumnRepository, *sqlx.DB, error)
	GetClientIndexRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientIndexRepository, *sqlx.DB, error)
	GetClientFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories.ClientFunctionRepository, *sqlx.DB, error)
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

func (s *ConnectionServiceImpl) GetClientTableRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientTableRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientTableRepo, err := repositories.NewClientTableRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientTableRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetClientColumnRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientColumnRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientColumnRepo, err := repositories.NewClientColumnRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientColumnRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetClientIndexRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientIndexRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientIndexRepo, err := repositories.NewClientIndexRepository(clientDatabaseConnection)
	if err != nil {
		return nil, nil, err
	}

	return clientIndexRepo, clientDatabaseConnection, nil
}

func (s *ConnectionServiceImpl) GetClientFunctionRepoByProjectUUID(projectUUID uuid.UUID, connection *sqlx.DB) (*repositories.ClientFunctionRepository, *sqlx.DB, error) {
	databaseName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return nil, nil, err
	}

	return s.getClientFunctionRepo(databaseName, connection)
}

func (s *ConnectionServiceImpl) getClientFunctionRepo(databaseName string, connection *sqlx.DB) (*repositories.ClientFunctionRepository, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientFunctionRepo, err := repositories.NewClientFunctionRepository(clientDatabaseConnection)
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
