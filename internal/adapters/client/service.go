package client

import (
	sqlxAdapter "fluxend/internal/adapters/sqlx"
	"fluxend/internal/database/repositories"
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/shared"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type ServiceImpl struct {
	databaseRepo shared.DatabaseService
}

func NewClientService(injector *do.Injector) (database.ConnectionService, error) {
	databaseRepo := do.MustInvoke[shared.DatabaseService](injector)

	return &ServiceImpl{
		databaseRepo: databaseRepo,
	}, nil
}

func (s *ServiceImpl) ConnectByDatabaseName(name string) (*sqlx.DB, error) {
	return s.databaseRepo.Connect(name)
}

func (s *ServiceImpl) GetDatabaseStatsRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	// Create a new injector with the client database connection
	clientInjector := do.New()

	// Register the raw sqlx.DB connection
	do.ProvideValue(clientInjector, clientDatabaseConnection)

	// Register the database adapter
	do.Provide(clientInjector, func(i *do.Injector) (shared.DB, error) {
		sqlxDB := do.MustInvoke[*sqlx.DB](i)
		return sqlxAdapter.NewAdapter(sqlxDB), nil
	})

	clientDatabaseStatsRepo, err := repositories.NewDatabaseStatsRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientDatabaseStatsRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetTableRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientInjector := s.createClientInjector(clientDatabaseConnection)

	clientTableRepo, err := repositories.NewTableRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientTableRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetFunctionRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientInjector := s.createClientInjector(clientDatabaseConnection)

	clientFunctionRepo, err := repositories.NewFunctionRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientFunctionRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetColumnRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientInjector := s.createClientInjector(clientDatabaseConnection)

	clientColumnRepo, err := repositories.NewColumnRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientColumnRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetIndexRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientInjector := s.createClientInjector(clientDatabaseConnection)

	clientIndexRepo, err := repositories.NewIndexRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientIndexRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) GetRowRepo(databaseName string, connection *sqlx.DB) (interface{}, *sqlx.DB, error) {
	clientDatabaseConnection, err := s.getOrCreateConnection(databaseName, connection)
	if err != nil {
		return nil, nil, err
	}

	clientInjector := s.createClientInjector(clientDatabaseConnection)

	clientRowRepo, err := repositories.NewRowRepository(clientInjector)
	if err != nil {
		return nil, nil, err
	}

	return clientRowRepo, clientDatabaseConnection, nil
}

func (s *ServiceImpl) getOrCreateConnection(databaseName string, connection *sqlx.DB) (*sqlx.DB, error) {
	if connection != nil {
		return connection, nil
	}

	return s.databaseRepo.Connect(databaseName)
}

// Helper method to create a dependency injector for client database connections
func (s *ServiceImpl) createClientInjector(clientConnection *sqlx.DB) *do.Injector {
	clientInjector := do.New()

	do.ProvideValue(clientInjector, clientConnection)

	do.Provide(clientInjector, func(i *do.Injector) (shared.DB, error) {
		sqlxDB := do.MustInvoke[*sqlx.DB](i)
		return sqlxAdapter.NewAdapter(sqlxDB), nil
	})

	return clientInjector
}
