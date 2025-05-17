package function

import (
	"fluxton/internal/adapters/connection"
	function2 "fluxton/internal/api/dto/database/function"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/models"
	"fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"

	"github.com/samber/do"
)

type FunctionService interface {
	List(schema string, projectUUID uuid.UUID, authUser models.AuthUser) ([]Function, error)
	GetByName(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (Function, error)
	Create(schema string, request *function2.CreateFunctionRequest, authUser models.AuthUser) (Function, error)
	Delete(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FunctionServiceImpl struct {
	connectService connection.ConnectionService
	projectPolicy  *project.ProjectPolicy
	databaseRepo   *repositories2.DatabaseRepository
	projectRepo    *repositories2.ProjectRepository
}

func NewFunctionService(injector *do.Injector) (FunctionService, error) {
	connectionService := do.MustInvoke[connection.ConnectionService](injector)
	policy := do.MustInvoke[*project.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories2.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &FunctionServiceImpl{
		connectService: connectionService,
		projectPolicy:  policy,
		databaseRepo:   databaseRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *FunctionServiceImpl) List(schema string, projectUUID uuid.UUID, authUser models.AuthUser) ([]Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.connectService.GetFunctionRepoByProjectUUID(projectUUID, nil)
	if err != nil {
		return []Function{}, err
	}
	defer connection.Close()

	return clientFunctionRepo.List(schema)
}

func (s *FunctionServiceImpl) GetByName(name, schema string, projectUUID uuid.UUID, authUser models.AuthUser) (Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.connectService.GetFunctionRepoByProjectUUID(projectUUID, nil)
	if err != nil {
		return Function{}, err
	}
	defer connection.Close()

	return clientFunctionRepo.GetByName(schema, name)
}

func (s *FunctionServiceImpl) Create(schema string, request *function2.CreateFunctionRequest, authUser models.AuthUser) (Function, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return Function{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.connectService.GetFunctionRepoByProjectUUID(request.ProjectUUID, nil)
	if err != nil {
		return Function{}, err
	}
	defer connection.Close()

	definitionQuery, err := s.buildDefinition(schema, request)
	if err != nil {
		return Function{}, err
	}

	err = clientFunctionRepo.Create(definitionQuery)
	if err != nil {
		return Function{}, err
	}

	return clientFunctionRepo.GetByName(schema, request.Name)
}

func (s *FunctionServiceImpl) Delete(schema, name string, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.connectService.GetFunctionRepoByProjectUUID(projectUUID, nil)
	if err != nil {
		return false, err
	}
	defer connection.Close()

	err = clientFunctionRepo.Delete(schema, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *FunctionServiceImpl) buildDefinition(schema string, request *function2.CreateFunctionRequest) (string, error) {
	var params []string
	for _, param := range request.Parameters {
		params = append(params, fmt.Sprintf("%s %s", pq.QuoteIdentifier(param.Name), pq.QuoteIdentifier(param.Type)))
	}
	paramList := strings.Join(params, ", ")

	sql := fmt.Sprintf(
		`CREATE OR REPLACE FUNCTION %s.%s(%s) RETURNS %s AS $$ %s; $$ LANGUAGE %s;`,
		schema,
		pq.QuoteIdentifier(request.Name),
		paramList,
		pq.QuoteIdentifier(request.ReturnType),
		pq.QuoteIdentifier(request.Definition),
		pq.QuoteIdentifier(request.Language),
	)

	return strings.ReplaceAll(sql, ";;", ";"), nil
}
