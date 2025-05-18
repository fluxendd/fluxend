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
	List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser auth.User) ([]Project, error)
	GetByUUID(projectUUID uuid.UUID, authUser auth.User) (Project, error)
	GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser auth.User) (string, error)
	Create(request *project.CreateRequest, authUser auth.User) (Project, error)
	Update(projectUUID uuid.UUID, authUser auth.User, request *project.UpdateRequest) (*Project, error)
	Delete(projectUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	projectPolicy    *Policy
	databaseRepo     client.Repository
	projectRepo      Repository
	postgrestService postgrest.Service
}

func NewProjectService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*Policy](injector)
	databaseRepo := do.MustInvoke[client.Repository](injector)
	projectRepo := do.MustInvoke[Repository](injector)
	postgrestService := do.MustInvoke[postgrest.Service](injector)

	return &ServiceImpl{
		projectPolicy:    policy,
		databaseRepo:     databaseRepo,
		projectRepo:      projectRepo,
		postgrestService: postgrestService,
	}, nil
}

func (s *ServiceImpl) List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser auth.User) ([]Project, error) {
	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Project{}, errors.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authUser.Uuid)
}

func (s *ServiceImpl) GetByUUID(projectUUID uuid.UUID, authUser auth.User) (Project, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return Project{}, err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return Project{}, errors.NewForbiddenError("project.error.viewForbidden")
	}

	return fetchedProject, nil
}

func (s *ServiceImpl) GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser auth.User) (string, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(fetchedProject.OrganizationUuid, authUser) {
		return "", errors.NewForbiddenError("project.error.viewForbidden")
	}

	return fetchedProject.DBName, nil
}

func (s *ServiceImpl) Create(request *project.CreateRequest, authUser auth.User) (Project, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationUUID, authUser) {
		return Project{}, errors.NewForbiddenError("project.error.createForbidden")
	}

	err := s.validateNameForDuplication(request.Name, request.OrganizationUUID)
	if err != nil {
		return Project{}, err
	}

	projectInput := Project{
		Name:             request.Name,
		OrganizationUuid: request.OrganizationUUID,
		DBName:           s.generateDBName(),
		DBPort:           s.generateDBPort(),
		CreatedBy:        authUser.Uuid,
		UpdatedBy:        authUser.Uuid,
	}

	_, err = s.projectRepo.Create(&projectInput)
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

func (s *ServiceImpl) Update(projectUUID uuid.UUID, authUser auth.User, request *project.UpdateRequest) (*Project, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return &Project{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	err = fetchedProject.PopulateModel(&fetchedProject, request)
	if err != nil {
		return nil, err
	}

	fetchedProject.UpdatedAt = time.Now()
	fetchedProject.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, fetchedProject.OrganizationUuid)
	if err != nil {
		return &Project{}, err
	}

	return s.projectRepo.Update(&fetchedProject)
}

func (s *ServiceImpl) Delete(projectUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedProject, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(fetchedProject.OrganizationUuid, authUser) {
		return false, errors.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.databaseRepo.DropIfExists(fetchedProject.DBName)
	if err != nil {
		return false, err
	}

	go s.postgrestService.RemoveContainer(fetchedProject.DBName)

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
