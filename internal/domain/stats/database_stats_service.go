package stats

import (
	"fluxton/internal/adapters/client"
	"fluxton/internal/domain/admin"
	"fluxton/internal/domain/database/stat"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/samber/do"
	"time"
)

type DatabaseStatsService interface {
	GetTotalDatabaseSize(databaseName string, authUser models.AuthUser) (string, error)
	GetTotalIndexSize(databaseName string, authUser models.AuthUser) (string, error)
	GetUnusedIndexes(databaseName string, authUser models.AuthUser) ([]stat.UnusedIndex, error)
	GetSlowQueries(databaseName string, authUser models.AuthUser) ([]stat.SlowQuery, error)
	GetIndexScansPerTable(databaseName string, authUser models.AuthUser) ([]stat.IndexScan, error)
	GetSizePerTable(databaseName string, authUser models.AuthUser) ([]stat.TableSize, error)
	GetRowCountPerTable(databaseName string, authUser models.AuthUser) ([]stat.TableRowCount, error)
	GetAll(databaseName string, authUser models.AuthUser) (stat.DatabaseStat, error)
}

type DatabaseStatsServiceImpl struct {
	connectionService client.Service
	adminPolicy       *admin.AdminPolicy
	databaseRepo      *client.Repository
}

func NewDatabaseStatsService(injector *do.Injector) (DatabaseStatsService, error) {
	connectionService := do.MustInvoke[client.Service](injector)
	policy := admin.NewAdminPolicy()
	databaseRepo := do.MustInvoke[*client.Repository](injector)

	return &DatabaseStatsServiceImpl{
		connectionService: connectionService,
		adminPolicy:       policy,
		databaseRepo:      databaseRepo,
	}, nil
}

func (s *DatabaseStatsServiceImpl) GetTotalDatabaseSize(databaseName string, authUser models.AuthUser) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalDatabaseSize()
}

func (s *DatabaseStatsServiceImpl) GetTotalIndexSize(databaseName string, authUser models.AuthUser) (string, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return "", errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalIndexSize()
}

func (s *DatabaseStatsServiceImpl) GetUnusedIndexes(databaseName string, authUser models.AuthUser) ([]stat.UnusedIndex, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []stat.UnusedIndex{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []stat.UnusedIndex{}, err
	}

	return dbStatsRepo.GetUnusedIndexes()
}

func (s *DatabaseStatsServiceImpl) GetSlowQueries(databaseName string, authUser models.AuthUser) ([]stat.SlowQuery, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []stat.SlowQuery{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []stat.SlowQuery{}, err
	}

	return dbStatsRepo.GetSlowQueries()
}

func (s *DatabaseStatsServiceImpl) GetIndexScansPerTable(databaseName string, authUser models.AuthUser) ([]stat.IndexScan, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []stat.IndexScan{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []stat.IndexScan{}, err
	}

	return dbStatsRepo.GetIndexScansPerTable()
}

func (s *DatabaseStatsServiceImpl) GetSizePerTable(databaseName string, authUser models.AuthUser) ([]stat.TableSize, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []stat.TableSize{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []stat.TableSize{}, err
	}

	return dbStatsRepo.GetSizePerTable()
}

func (s *DatabaseStatsServiceImpl) GetRowCountPerTable(databaseName string, authUser models.AuthUser) ([]stat.TableRowCount, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return []stat.TableRowCount{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	dbStatsRepo, _, err := s.connectionService.GetDatabaseStatsRepo(databaseName, nil)
	if err != nil {
		return []stat.TableRowCount{}, err
	}

	return dbStatsRepo.GetRowCountPerTable()
}

func (s *DatabaseStatsServiceImpl) GetAll(databaseName string, authUser models.AuthUser) (stat.DatabaseStat, error) {
	if !s.adminPolicy.CanAccess(authUser) {
		return stat.DatabaseStat{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	totalDatabaseSize, err := s.GetTotalDatabaseSize(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	totalIndexSize, err := s.GetTotalIndexSize(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	unusedIndexes, err := s.GetUnusedIndexes(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	tableCounts, err := s.GetRowCountPerTable(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	tableSizes, err := s.GetSizePerTable(databaseName, authUser)
	if err != nil {
		return stat.DatabaseStat{}, err
	}

	return stat.DatabaseStat{
		DatabaseName: databaseName,
		TotalSize:    totalDatabaseSize,
		IndexSize:    totalIndexSize,
		UnusedIndex:  unusedIndexes,
		TableCount:   tableCounts,
		TableSize:    tableSizes,
		CreatedAt:    time.Now(),
	}, nil
}
