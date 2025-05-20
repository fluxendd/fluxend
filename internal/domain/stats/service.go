package stats

import (
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/database"
	"fluxton/internal/domain/shared"
	"fluxton/pkg/errors"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type Service interface {
	GetTotalDatabaseSize(databaseName string, authUser auth.User) (string, error)
	GetTotalIndexSize(databaseName string, authUser auth.User) (string, error)
	GetUnusedIndexes(databaseName string, authUser auth.User) ([]UnusedIndex, error)
	GetSlowQueries(databaseName string, authUser auth.User) ([]SlowQuery, error)
	GetIndexScansPerTable(databaseName string, authUser auth.User) ([]IndexScan, error)
	GetSizePerTable(databaseName string, authUser auth.User) ([]TableSize, error)
	GetRowCountPerTable(databaseName string, authUser auth.User) ([]TableRowCount, error)
	GetAll(databaseName string, authUser auth.User) (Stat, error)
}

type ServiceImpl struct {
	connectionService database.ConnectionService
	adminPolicy       *admin.Policy
	databaseRepo      *shared.DatabaseService
}

func NewDatabaseStatsService(injector *do.Injector) (Service, error) {
	connectionService := do.MustInvoke[database.ConnectionService](injector)
	policy := admin.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*shared.DatabaseService](injector)

	return &ServiceImpl{
		connectionService: connectionService,
		adminPolicy:       policy,
		databaseRepo:      databaseRepo,
	}, nil
}

func (s *ServiceImpl) GetTotalDatabaseSize(databaseName string, authUser auth.User) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalDatabaseSize()
}

func (s *ServiceImpl) GetTotalIndexSize(databaseName string, authUser auth.User) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalIndexSize()
}

func (s *ServiceImpl) GetUnusedIndexes(databaseName string, authUser auth.User) ([]UnusedIndex, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []UnusedIndex{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []UnusedIndex{}, err
	}

	return dbStatsRepo.GetUnusedIndexes()
}

func (s *ServiceImpl) GetSlowQueries(databaseName string, authUser auth.User) ([]SlowQuery, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []SlowQuery{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []SlowQuery{}, err
	}

	return dbStatsRepo.GetSlowQueries()
}

func (s *ServiceImpl) GetIndexScansPerTable(databaseName string, authUser auth.User) ([]IndexScan, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []IndexScan{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []IndexScan{}, err
	}

	return dbStatsRepo.GetIndexScansPerTable()
}

func (s *ServiceImpl) GetSizePerTable(databaseName string, authUser auth.User) ([]TableSize, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []TableSize{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []TableSize{}, err
	}

	return dbStatsRepo.GetSizePerTable()
}

func (s *ServiceImpl) GetRowCountPerTable(databaseName string, authUser auth.User) ([]TableRowCount, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []TableRowCount{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []TableRowCount{}, err
	}

	return dbStatsRepo.GetRowCountPerTable()
}

func (s *ServiceImpl) GetAll(databaseName string, authUser auth.User) (Stat, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return Stat{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	totalDatabaseSize, err := s.GetTotalDatabaseSize(databaseName, authUser)
	if err != nil {
		return Stat{}, err
	}

	totalIndexSize, err := s.GetTotalIndexSize(databaseName, authUser)
	if err != nil {
		return Stat{}, err
	}

	unusedIndexes, err := s.GetUnusedIndexes(databaseName, authUser)
	if err != nil {
		return Stat{}, err
	}

	tableCounts, err := s.GetRowCountPerTable(databaseName, authUser)
	if err != nil {
		return Stat{}, err
	}

	tableSizes, err := s.GetSizePerTable(databaseName, authUser)
	if err != nil {
		return Stat{}, err
	}

	return Stat{
		DatabaseName: databaseName,
		TotalSize:    totalDatabaseSize,
		IndexSize:    totalIndexSize,
		UnusedIndex:  unusedIndexes,
		TableCount:   tableCounts,
		TableSize:    tableSizes,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *ServiceImpl) getClientStatsRepo(dbName string) (StatRepository, *sqlx.DB, error) {
	repo, connection, err := s.connectionService.GetDatabaseStatsRepo(dbName, nil)
	if err != nil {
		return nil, nil, err
	}

	clientRepo, ok := repo.(StatRepository)
	if !ok {
		connection.Close()

		return nil, nil, errors.NewUnprocessableError("clientStatsRepo is invalid")
	}

	return clientRepo, connection, nil
}
