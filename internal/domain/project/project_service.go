package project

import (
	"fluxton/internal/adapters/postgrest"
	"fluxton/internal/api/dto"
	project2 "fluxton/internal/api/dto/project"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"math/rand"
	"strings"
	"time"
)

type ProjectService interface {
	List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser models.AuthUser) ([]Project, error)
	GetByUUID(projectUUID uuid.UUID, authUser models.AuthUser) (Project, error)
	GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser models.AuthUser) (string, error)
	Create(request *project2.CreateRequest, authUser models.AuthUser) (Project, error)
	Update(projectUUID uuid.UUID, authUser models.AuthUser, request *project2.UpdateRequest) (*Project, error)
	Delete(projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ProjectServiceImpl struct {
	projectPolicy    *ProjectPolicy
	databaseRepo     *repositories2.DatabaseRepository
	projectRepo      *repositories2.ProjectRepository
	postgrestService postgrest.PostgrestService
}

func NewProjectService(injector *do.Injector) (ProjectService, error) {
	policy := do.MustInvoke[*ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories2.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)
	postgrestService := do.MustInvoke[postgrest.PostgrestService](injector)

	return &ProjectServiceImpl{
		projectPolicy:    policy,
		databaseRepo:     databaseRepo,
		projectRepo:      projectRepo,
		postgrestService: postgrestService,
	}, nil
}

func (s *ProjectServiceImpl) List(paginationParams dto.PaginationParams, organizationUUID uuid.UUID, authUser models.AuthUser) ([]Project, error) {
	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Project{}, errors.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authUser.Uuid)
}

func (s *ProjectServiceImpl) GetByUUID(projectUUID uuid.UUID, authUser models.AuthUser) (Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return project.Project{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return project.Project{}, errors.NewForbiddenError("project.error.viewForbidden")
	}

	return project, nil
}

func (s *ProjectServiceImpl) GetDatabaseNameByUUID(projectUUID uuid.UUID, authUser models.AuthUser) (string, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return "", err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationUuid, authUser) {
		return "", errors.NewForbiddenError("project.error.viewForbidden")
	}

	return project.DBName, nil
}

func (s *ProjectServiceImpl) Create(request *project2.CreateRequest, authUser models.AuthUser) (Project, error) {
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
		return project.Project{}, err
	}

	err = s.databaseRepo.Create(project.DBName, uuid.NullUUID{UUID: authUser.Uuid, Valid: true})
	if err != nil {
		// TODO: handle better
		s.projectRepo.Delete(project.Uuid)

		return project.Project{}, err
	}

	go s.postgrestService.StartContainer(project.DBName)

	return project, nil
}

func (s *ProjectServiceImpl) Update(projectUUID uuid.UUID, authUser models.AuthUser, request *project2.UpdateRequest) (*Project, error) {
	project, err := s.projectRepo.GetByUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationUuid, authUser) {
		return &project.Project{}, errors.NewForbiddenError("project.error.updateForbidden")
	}

	err = project.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	project.UpdatedAt = time.Now()
	project.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, project.OrganizationUuid)
	if err != nil {
		return &project.Project{}, err
	}

	return s.projectRepo.Update(&project)
}

func (s *ProjectServiceImpl) Delete(projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
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

func (s *ProjectServiceImpl) generateDBName() string {
	return "udb_" + strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "")
}

func (s *ProjectServiceImpl) generateDBPort() int {
	return rand.Intn(65535-5000+1) + 5000
}

func (s *ProjectServiceImpl) validateNameForDuplication(name string, organizationUUID uuid.UUID) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
