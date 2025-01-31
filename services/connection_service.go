package services

import (
	"fluxton/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type ConnectionService interface {
	ConnectByDatabaseName(name string) (*sqlx.DB, error)
	ConnectByProjectID(projectID uint) (*sqlx.DB, error)
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

func (s *ConnectionServiceImpl) ConnectByProjectID(projectID uint) (*sqlx.DB, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	return s.databaseRepo.Connect(project.DBName)
}
