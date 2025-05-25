package database

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/shared"
	"fluxend/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"

	"github.com/samber/do"
)

type FunctionService interface {
	List(schema string, projectUUID uuid.UUID, authUser auth.User) ([]Function, error)
	GetByName(name, schema string, projectUUID uuid.UUID, authUser auth.User) (Function, error)
	Create(schema string, request CreateFunctionInput, authUser auth.User) (Function, error)
	Delete(name, schema string, projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type FunctionServiceImpl struct {
	connectionService ConnectionService
	projectPolicy     *project.Policy
	databaseRepo      shared.DatabaseService
	projectRepo       project.Repository
}

func NewFunctionService(injector *do.Injector) (FunctionService, error) {
	connectionService := do.MustInvoke[ConnectionService](injector)
	policy := do.MustInvoke[*project.Policy](injector)
	databaseRepo := do.MustInvoke[shared.DatabaseService](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &FunctionServiceImpl{
		connectionService: connectionService,
		projectPolicy:     policy,
		databaseRepo:      databaseRepo,
		projectRepo:       projectRepo,
	}, nil
}

func (s *FunctionServiceImpl) List(schema string, projectUUID uuid.UUID, authUser auth.User) ([]Function, error) {
	dbName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return []Function{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.getClientFunctionRepo(dbName)
	if err != nil {
		return []Function{}, err
	}
	defer connection.Close()

	return clientFunctionRepo.List(schema)
}

func (s *FunctionServiceImpl) GetByName(name, schema string, projectUUID uuid.UUID, authUser auth.User) (Function, error) {
	dbName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return Function{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return Function{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.getClientFunctionRepo(dbName)
	if err != nil {
		return Function{}, err
	}
	defer connection.Close()

	return clientFunctionRepo.GetByName(schema, name)
}

func (s *FunctionServiceImpl) Create(schema string, request CreateFunctionInput, authUser auth.User) (Function, error) {
	dbName, err := s.projectRepo.GetDatabaseNameByUUID(request.ProjectUUID)
	if err != nil {
		return Function{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return Function{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return Function{}, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.getClientFunctionRepo(dbName)
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

func (s *FunctionServiceImpl) Delete(schema, name string, projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	dbName, err := s.projectRepo.GetDatabaseNameByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("function.error.listForbidden")
	}

	clientFunctionRepo, connection, err := s.getClientFunctionRepo(dbName)
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

func (s *FunctionServiceImpl) buildDefinition(schema string, request CreateFunctionInput) (string, error) {
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

func (s *FunctionServiceImpl) getClientFunctionRepo(dbName string) (FunctionRepository, *sqlx.DB, error) {
	repo, connection, err := s.connectionService.GetFunctionRepo(dbName, nil)
	if err != nil {
		return nil, nil, err
	}

	clientRepo, ok := repo.(FunctionRepository)
	if !ok {
		connection.Close()

		return nil, nil, errors.NewUnprocessableError("clientFunctionRepo is invalid")
	}

	return clientRepo, connection, nil
}
