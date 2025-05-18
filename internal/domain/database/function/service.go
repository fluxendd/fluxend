package function

import (
	"fluxton/internal/adapters/client"
	"fluxton/internal/api/dto/database/function"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"

	"github.com/samber/do"
)

type Service interface {
	List(schema string, projectUUID uuid.UUID, authUser auth.User) ([]Function, error)
	GetByName(name, schema string, projectUUID uuid.UUID, authUser auth.User) (Function, error)
	Create(schema string, request *function.CreateFunctionRequest, authUser auth.User) (Function, error)
	Delete(name, schema string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	connectService client.Service
	projectPolicy  *project.Policy
	databaseRepo   client.Repository
	projectRepo    project.Repository
}

func NewFunctionService(injector *do.Injector) (Service, error) {
	connectionService := do.MustInvoke[client.Service](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	databaseRepo := do.MustInvoke[client.Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		connectService: connectionService,
		projectPolicy:  policy,
		databaseRepo:   databaseRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *ServiceImpl) List(schema string, projectUUID uuid.UUID, authUser auth.User) ([]Function, error) {
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

func (s *ServiceImpl) GetByName(name, schema string, projectUUID uuid.UUID, authUser auth.User) (Function, error) {
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

func (s *ServiceImpl) Create(schema string, request *function.CreateFunctionRequest, authUser auth.User) (Function, error) {
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

func (s *ServiceImpl) Delete(schema, name string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
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

func (s *ServiceImpl) buildDefinition(schema string, request *function.CreateFunctionRequest) (string, error) {
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
