package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fmt"
	"github.com/google/uuid"
	"strings"

	"github.com/samber/do"
)

type FunctionService interface {
	List(schema string, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Function, error)
	GetByName(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (models.Function, error)
	Create(request *requests.CreateFunctionRequest, projectUUID uuid.UUID, authUser models.AuthUser) (models.Function, error)
	Delete(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FunctionServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	databaseRepo  *repositories.DatabaseRepository
	projectRepo   *repositories.ProjectRepository
}

func NewFunctionService(injector *do.Injector) (FunctionService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FunctionServiceImpl{
		projectPolicy: policy,
		databaseRepo:  databaseRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FunctionServiceImpl) List(schema string, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Function{}, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.getClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return []models.Function{}, err
	}

	return clientFunctionRepo.List(schema)
}

func (s *FunctionServiceImpl) GetByName(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (models.Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Function{}, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.getClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return models.Function{}, err
	}

	return clientFunctionRepo.GetByName(schema, name)
}

func (s *FunctionServiceImpl) Create(request *requests.CreateFunctionRequest, projectUUID uuid.UUID, authUser models.AuthUser) (models.Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.Function{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.Function{}, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.getClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return models.Function{}, err
	}

	definitionQuery, err := s.buildDefinition(request)
	if err != nil {
		return models.Function{}, err
	}

	err = clientFunctionRepo.Create(definitionQuery)
	if err != nil {
		return models.Function{}, err
	}

	return clientFunctionRepo.GetByName(request.Schema, request.Name)
}

func (s *FunctionServiceImpl) Delete(schema, name string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.getClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	err = clientFunctionRepo.Delete(schema, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *FunctionServiceImpl) buildDefinition(request *requests.CreateFunctionRequest) (string, error) {
	var params []string
	for _, param := range request.Parameters {
		params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
	}
	paramList := strings.Join(params, ", ")

	sql := fmt.Sprintf(
		`CREATE OR REPLACE FUNCTION %s.%s(%s) RETURNS %s AS $$ BEGIN %s END; $$ LANGUAGE %s;`,
		request.Schema,
		request.Name,
		paramList,
		request.ReturnType,
		request.Definition,
		request.Language,
	)

	return sql, nil
}

func (s *FunctionServiceImpl) getClientFunctionRepoByProjectUUID(projectUUID uuid.UUID) (*repositories.ClientFunctionRepository, error) {
	databaseName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	return s.getClientFunctionRepo(databaseName)
}

func (s *FunctionServiceImpl) getClientFunctionRepo(databaseName string) (*repositories.ClientFunctionRepository, error) {
	clientDatabaseConnection, err := s.databaseRepo.Connect(databaseName)
	if err != nil {
		return nil, err
	}

	clientFunctionRepo, err := repositories.NewClientFunctionRepository(clientDatabaseConnection)
	if err != nil {
		return nil, err
	}

	return clientFunctionRepo, nil
}
