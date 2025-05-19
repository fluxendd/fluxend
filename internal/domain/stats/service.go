package stats

import (
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/database"
	"fluxton/internal/domain/shared"
	"fluxton/pkg/errors"
	"github.com/samber/do"
	"time"
)

type Service interface {
	GetTotalDatabaseSize(databaseName string, authUser auth.User) (string, error)
	GetTotalIndexSize(databaseName string, authUser auth.User) (string, error)
	GetUnusedIndexes(databaseName string, authUser auth.User) ([]database.UnusedIndex, error)
	GetSlowQueries(databaseName string, authUser auth.User) ([]database.SlowQuery, error)
	GetIndexScansPerTable(databaseName string, authUser auth.User) ([]database.IndexScan, error)
	GetSizePerTable(databaseName string, authUser auth.User) ([]database.TableSize, error)
	GetRowCountPerTable(databaseName string, authUser auth.User) ([]database.TableRowCount, error)
	GetAll(databaseName string, authUser auth.User) (database.Stat, error)
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

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalDatabaseSize()
}

func (s *ServiceImpl) GetTotalIndexSize(databaseName string, authUser auth.User) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalIndexSize()
}

func (s *ServiceImpl) GetUnusedIndexes(databaseName string, authUser auth.User) ([]database.UnusedIndex, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []database.UnusedIndex{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []database.UnusedIndex{}, err
	}

	return dbStatsRepo.GetUnusedIndexes()
}

func (s *ServiceImpl) GetSlowQueries(databaseName string, authUser auth.User) ([]database.SlowQuery, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []database.SlowQuery{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []database.SlowQuery{}, err
	}

	return dbStatsRepo.GetSlowQueries()
}

func (s *ServiceImpl) GetIndexScansPerTable(databaseName string, authUser auth.User) ([]database.IndexScan, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []database.IndexScan{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []database.IndexScan{}, err
	}

	return dbStatsRepo.GetIndexScansPerTable()
}

func (s *ServiceImpl) GetSizePerTable(databaseName string, authUser auth.User) ([]database.TableSize, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []database.TableSize{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []database.TableSize{}, err
	}

	return dbStatsRepo.GetSizePerTable()
}

func (s *ServiceImpl) GetRowCountPerTable(databaseName string, authUser auth.User) ([]database.TableRowCount, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []database.TableRowCount{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []database.TableRowCount{}, err
	}

	return dbStatsRepo.GetRowCountPerTable()
}

func (s *ServiceImpl) GetAll(databaseName string, authUser auth.User) (database.Stat, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return database.Stat{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	totalDatabaseSize, err := s.GetTotalDatabaseSize(databaseName, authUser)
	if err != nil {
		return database.Stat{}, err
	}

	totalIndexSize, err := s.GetTotalIndexSize(databaseName, authUser)
	if err != nil {
		return database.Stat{}, err
	}

	unusedIndexes, err := s.GetUnusedIndexes(databaseName, authUser)
	if err != nil {
		return database.Stat{}, err
	}

	tableCounts, err := s.GetRowCountPerTable(databaseName, authUser)
	if err != nil {
		return database.Stat{}, err
	}

	tableSizes, err := s.GetSizePerTable(databaseName, authUser)
	if err != nil {
		return database.Stat{}, err
	}

	return database.Stat{
		DatabaseName: databaseName,
		TotalSize:    totalDatabaseSize,
		IndexSize:    totalIndexSize,
		UnusedIndex:  unusedIndexes,
		TableCount:   tableCounts,
		TableSize:    tableSizes,
		CreatedAt:    time.Now(),
	}, nil
}
