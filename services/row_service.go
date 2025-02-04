package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
	"strings"
)

type RowService interface {
	List(paginationParams utils.PaginationParams, tableName string, organizationID, projectID, authenticatedUserID uuid.UUID) ([]models.Row, error)
	Create(request *requests.RowCreateRequest, projectID uuid.UUID, tableName string, authenticatedUser models.AuthenticatedUser) (models.Row, error)
}

type RowServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	databaseRepo  *repositories.DatabaseRepository
	projectRepo   *repositories.ProjectRepository
	coreTableRepo *repositories.CoreTableRepository
}

func NewRowService(injector *do.Injector) (RowService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &RowServiceImpl{
		projectPolicy: policy,
		databaseRepo:  databaseRepo,
		projectRepo:   projectRepo,
		coreTableRepo: coreTableRepo,
	}, nil
}

func (s *RowServiceImpl) List(paginationParams utils.PaginationParams, tableName string, organizationID, projectID, authenticatedUserID uuid.UUID) ([]models.Row, error) {
	if !s.projectPolicy.CanList(organizationID, authenticatedUserID) {
		return []models.Row{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return []models.Row{}, err
	}

	clientRowRepo, err := s.getClientRowRepo(project.DBName)
	if err != nil {
		return []models.Row{}, err
	}

	return clientRowRepo.List(tableName, paginationParams)
}

func (s *RowServiceImpl) Create(request *requests.RowCreateRequest, projectID uuid.UUID, tableName string, authenticatedUser models.AuthenticatedUser) (models.Row, error) {
	table, err := s.coreTableRepo.GetByName(tableName)
	if err != nil {
		return models.Row{}, err
	}

	err = s.validateColumns(request, table)
	if err != nil {
		return models.Row{}, err
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Row{}, err
	}

	clientRowRepo, err := s.getClientRowRepo(project.DBName)
	if err != nil {
		return models.Row{}, err
	}

	insertedID, err := clientRowRepo.Create(table.Name, request.Fields)
	if err != nil {
		return models.Row{}, err
	}

	row, err := clientRowRepo.GetByID(table.Name, insertedID)
	if err != nil {
		return models.Row{}, err
	}

	return row, nil
}

func (s *RowServiceImpl) validateColumns(request *requests.RowCreateRequest, table models.Table) error {
	for f := range request.Fields {
		fieldName := strings.TrimSpace(f)
		columnFound := false

		for _, column := range table.Columns {
			if column.Name == fieldName {
				columnFound = true
				break
			}
		}

		if !columnFound {
			return errs.NewUnprocessableError(fmt.Sprintf("row.error.columnNotFound: %s", fieldName))
		}
	}

	return nil
}

func (s *RowServiceImpl) getClientRowRepo(databaseName string) (*repositories.ClientRowRepository, error) {
	clientDatabaseConnection, err := s.databaseRepo.Connect(databaseName)
	if err != nil {
		return nil, err
	}

	clientRowRepo, err := repositories.NewClientRowRepository(clientDatabaseConnection)
	if err != nil {
		return nil, err
	}

	return clientRowRepo, nil
}
