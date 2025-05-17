package project

import (
	"fluxton/internal/adapters/client"
	"fluxton/internal/adapters/postgrest"
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/project"
	"fluxton/internal/domain/auth"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"math/rand"
	"strings"
	"time"
)

type Service interface {
	List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser auth.AuthUser) ([]Project, error)
	GetByUUID(projectUUID uuid.UUID, authUser auth.AuthUser) (Project, error)
	GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser auth.AuthUser) (string, error)
	Create(request *project.CreateRequest, authUser auth.AuthUser) (Project, error)
	Update(projectUUID uuid.UUID, authUser auth.AuthUser, request *project.UpdateRequest) (*Project, error)
	Delete(projectUUID uuid.UUID, authUser auth.AuthUser) (bool, error)
}

type ServiceImpl struct {
	projectPolicy    *Policy
	databaseRepo     *client.Repository
	projectRepo      *Repository
	postgrestService postgrest.Service
}

func NewProjectService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*Policy](injector)
	databaseRepo := do.MustInvoke[*client.Repository](injector)
	projectRepo := do.MustInvoke[*Repository](injector)
	postgrestService := do.MustInvoke[postgrest.Service](injector)

	return &ServiceImpl{
		projectPolicy:    policy,
		databaseRepo:     databaseRepo,
		projectRepo:      projectRepo,
		postgrestService: postgrestService,
	}, nil
}

func (s *ServiceImpl) List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser auth.AuthUser) ([]Project, error) {
	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Project{}, errors.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authUser.Uuid)
}

func (s *ServiceImpl) GetByUUID(projectUUID uuid.UUID, authUser auth.AuthUser) (Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Project{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return Project{}, errors.NewForbiddenError("project.error.viewForbidden")
	}

	return project, nil
}

func (s *ServiceImpl) GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser auth.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return "", errors.NewForbiddenError("project.error.viewForbidden")
	}

	return project.DBName, nil
}

func (s *ServiceImpl) Create(request *project.CreateRequest, authUser auth.AuthUser) (Project, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationUUID, authUser) {
		return Project{}, errors.NewForbiddenError("project.error.createForbidden")
	}

	err := s.validateNameForDuplication(request.Name, request.OrganizationUUID)
	if err != nil {
		return Project{}, err
	}

	project := Project{
		Name:             request.Name,
		OrganizationUuid: request.OrganizationUUID,
		DBName:           s.generateDBName(),
		DBPort:           s.generateDBPort(),
		CreatedBy:        authUser.Uuid,
		UpdatedBy:        authUser.Uuid,
	}

	_, err = s.projectRepo.Create(&project)
	if err != nil {
		return Project{}, err
	}

	err = s.databaseRepo.Create(project.DBName, uuid.NullUUID{UUID: authUser.Uuid, Valid: true})
	if err != nil {
		// TODO: handle better
		s.projectRepo.Delete(project.Uuid)

		return Project{}, err
	}

	go s.postgrestService.StartContainer(project.DBName)

	return project, nil
}

func (s *ServiceImpl) Update(projectUUID uuid.UUID, authUser auth.AuthUser, request *project.UpdateRequest) (*Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &Project{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	err = project.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	project.UpdatedAt = time.Now()
	project.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, project.OrganizationUuid)
	if err != nil {
		return &Project{}, err
	}

	return s.projectRepo.Update(&project)
}

func (s *ServiceImpl) Delete(projectUUID uuid.UUID, authUser auth.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return false, errors.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.databaseRepo.DropIfExists(project.DBName)
	if err != nil {
		return false, err
	}

	go s.postgrestService.RemoveContainer(project.DBName)

	return s.projectRepo.Delete(projectUUID)
}

func (s *ServiceImpl) generateDBName() string {
	return "udb_" + strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "")
}

func (s *ServiceImpl) generateDBPort() int {
	return rand.Intn(65535-5000+1) + 5000
}

func (s *ServiceImpl) validateNameForDuplication(name string, organizationUUID uuid.UUID) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
