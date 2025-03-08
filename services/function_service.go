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
	Create(schema string, request *requests.CreateFunctionRequest, authUser models.AuthUser) (models.Function, error)
	Delete(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FunctionServiceImpl struct {
	connectService ConnectionService
	projectPolicy  *policies.ProjectPolicy
	databaseRepo   *repositories.DatabaseRepository
	projectRepo    *repositories.ProjectRepository
}

func NewFunctionService(injector *do.Injector) (FunctionService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FunctionServiceImpl{
		connectService: connectionService,
		projectPolicy:  policy,
		databaseRepo:   databaseRepo,
		projectRepo:    projectRepo,
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

	clientFunctionRepo, err := s.connectService.GetClientFunctionRepoByProjectUUID(projectUUID)
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

	clientFunctionRepo, err := s.connectService.GetClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return models.Function{}, err
	}

	return clientFunctionRepo.GetByName(schema, name)
}

func (s *FunctionServiceImpl) Create(schema string, request *requests.CreateFunctionRequest, authUser models.AuthUser) (models.Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return models.Function{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.Function{}, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.connectService.GetClientFunctionRepoByProjectUUID(request.ProjectUUID)
	if err != nil {
		return models.Function{}, err
	}

	definitionQuery, err := s.buildDefinition(schema, request)
	if err != nil {
		return models.Function{}, err
	}

	err = clientFunctionRepo.Create(definitionQuery)
	if err != nil {
		return models.Function{}, err
	}

	return clientFunctionRepo.GetByName(schema, request.Name)
}

func (s *FunctionServiceImpl) Delete(schema, name string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, err := s.connectService.GetClientFunctionRepoByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	err = clientFunctionRepo.Delete(schema, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *FunctionServiceImpl) buildDefinition(schema string, request *requests.CreateFunctionRequest) (string, error) {
	var params []string
	for _, param := range request.Parameters {
		params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
	}
	paramList := strings.Join(params, ", ")

	sql := fmt.Sprintf(
		`CREATE OR REPLACE FUNCTION %s.%s(%s) RETURNS %s AS $$ %s; $$ LANGUAGE %s;`,
		schema,
		request.Name,
		paramList,
		request.ReturnType,
		request.Definition,
		request.Language,
	)

	return strings.ReplaceAll(sql, ";;", ";"), nil
}
