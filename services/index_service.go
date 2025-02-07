package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type IndexService interface {
	List(tableID, projectID uuid.UUID, authUser models.AuthUser) ([]string, error)
	GetByName(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (string, error)
	Create(projectID, tableID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (models.Table, error)
	Delete(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type IndexServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *policies.ProjectPolicy
	projectRepo       *repositories.ProjectRepository
	coreTableRepo     *repositories.CoreTableRepository
}

func NewIndexService(injector *do.Injector) (IndexService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	coreTableRepo := do.MustInvoke[*repositories.CoreTableRepository](injector)

	return &IndexServiceImpl{
		projectPolicy:     policy,
		connectionService: connectionService,
		projectRepo:       projectRepo,
		coreTableRepo:     coreTableRepo,
	}, nil
}

func (s *IndexServiceImpl) List(tableID, projectID uuid.UUID, authUser models.AuthUser) ([]string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationID, authUser) {
		return nil, errs.NewForbiddenError("project.error.readForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return nil, err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return nil, err
	}

	return clientIndexRepo.List(table.Name)
}

func (s *IndexServiceImpl) GetByName(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationID, authUser) {
		return "", errs.NewForbiddenError("project.error.readForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return "", err
	}

	clientIndexRepo, err := s.connectionService.GetClientIndexRepo(project.DBName)
	if err != nil {
		return "", err
	}

	return clientIndexRepo.GetByName(table.Name, indexName)
}

func (s *IndexServiceImpl) Create(projectID, tableID uuid.UUID, request *requests.IndexCreateRequest, authUser models.AuthUser) (models.Table, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Table{}, err
	}

	if !s.projectPolicy.CanCreate(request.OrganizationID, authUser) {
		return models.Table{}, errs.NewForbiddenError("table.error.createForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return models.Table{}, err
	}

	err = s.validateNameForDuplication(request.Column.Name, tableID)
	if err != nil {
		return models.Table{}, err
	}

	table.Columns = append(table.Columns, request.Column)
	table.UpdatedBy = authUser.ID

	_, err = s.coreTableRepo.Update(&table)
	if err != nil {
		return models.Table{}, err
	}

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return models.Table{}, err
	}

	err = clientColumnRepo.Create(table.Name, request.Column)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (s *IndexServiceImpl) Delete(indexName string, tableID, projectID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationID, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	table, err := s.coreTableRepo.GetByID(tableID)
	if err != nil {
		return false, err
	}

	for i, column := range table.Columns {
		if column.Name == columnName {
			// Remove column from slice
			table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
			break
		}
	}

	table.UpdatedBy = authUser.ID

	clientColumnRepo, err := s.connectionService.GetClientColumnRepo(project.DBName)
	if err != nil {
		return false, err
	}

	err = clientColumnRepo.Drop(table.Name, columnName)
	if err != nil {
		return false, err
	}

	_, err = s.coreTableRepo.Update(&table)

	return err == nil, err
}

func (s *IndexServiceImpl) validateNameForDuplication(name string, tableID uuid.UUID) error {
	exists, err := s.coreTableRepo.HasColumn(name, tableID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("table.error.duplicateName")
	}

	return nil
}
