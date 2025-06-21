package stats

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/shared"
	"fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type Service interface {
	GetAll(projectUUID uuid.UUID, authUser auth.User) (Stat, error)
}

type ServiceImpl struct {
	connectionService database.ConnectionService
	projectPolicy     *project.Policy
	databaseRepo      shared.DatabaseService
	projectRepo       project.Repository
}

func NewDatabaseStatsService(injector *do.Injector) (Service, error) {
	connectionService := do.MustInvoke[database.ConnectionService](injector)
	databaseRepo := do.MustInvoke[shared.DatabaseService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		connectionService: connectionService,
		projectPolicy:     policy,
		databaseRepo:      databaseRepo,
		projectRepo:       projectRepo,
	}, nil
}

func (s *ServiceImpl) getTotalDatabaseSize(databaseName string) (string, error) {
	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalDatabaseSize()
}

func (s *ServiceImpl) getTotalIndexSize(databaseName string) (string, error) {
	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return "", err
	}

	return dbStatsRepo.GetTotalIndexSize()
}

func (s *ServiceImpl) getUnusedIndexes(databaseName string) ([]UnusedIndex, error) {
	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []UnusedIndex{}, err
	}

	return dbStatsRepo.GetUnusedIndexes()
}

func (s *ServiceImpl) getSizePerTable(databaseName string) ([]TableSize, error) {
	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []TableSize{}, err
	}

	return dbStatsRepo.GetSizePerTable()
}

func (s *ServiceImpl) getRowCountPerTable(databaseName string) ([]TableRowCount, error) {
	dbStatsRepo, _, err := s.getClientStatsRepo(databaseName)
	if err != nil {
		return []TableRowCount{}, err
	}

	return dbStatsRepo.GetRowCountPerTable()
}

func (s *ServiceImpl) GetAll(projectUUID uuid.UUID, authUser auth.User) (Stat, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Stat{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return Stat{}, errors.NewForbiddenError("database_stats.error.forbidden")
	}

	databaseName := fetchedProject.DBName
	totalDatabaseSize, err := s.getTotalDatabaseSize(databaseName)
	if err != nil {
		return Stat{}, err
	}

	totalIndexSize, err := s.getTotalIndexSize(databaseName)
	if err != nil {
		return Stat{}, err
	}

	unusedIndexes, err := s.getUnusedIndexes(databaseName)
	if err != nil {
		return Stat{}, err
	}

	tableCounts, err := s.getRowCountPerTable(databaseName)
	if err != nil {
		return Stat{}, err
	}

	tableSizes, err := s.getSizePerTable(databaseName)
	if err != nil {
		return Stat{}, err
	}

	return Stat{
		DatabaseName: fetchedProject.DBName,
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
