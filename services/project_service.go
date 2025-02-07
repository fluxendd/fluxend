package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"strings"
	"time"
)

type ProjectService interface {
	List(paginationParams utils.PaginationParams, organizationID uuid.UUID, authUser models.AuthUser) ([]models.Project, error)
	GetByID(projectID uuid.UUID, authUser models.AuthUser) (models.Project, error)
	Create(request *requests.ProjectCreateRequest, authUser models.AuthUser) (models.Project, error)
	Update(projectID uuid.UUID, authUser models.AuthUser, request *requests.ProjectUpdateRequest) (*models.Project, error)
	Delete(projectID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type ProjectServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	databaseRepo  *repositories.DatabaseRepository
	projectRepo   *repositories.ProjectRepository
}

func NewProjectService(injector *do.Injector) (ProjectService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	databaseRepo := do.MustInvoke[*repositories.DatabaseRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &ProjectServiceImpl{
		projectPolicy: policy,
		databaseRepo:  databaseRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *ProjectServiceImpl) List(paginationParams utils.PaginationParams, organizationID uuid.UUID, authUser models.AuthUser) ([]models.Project, error) {
	if !s.projectPolicy.CanAccess(organizationID, authUser) {
		return []models.Project{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authUser.ID)
}

func (s *ProjectServiceImpl) GetByID(projectID uuid.UUID, authUser models.AuthUser) (models.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return models.Project{}, err
	}

	if !s.projectPolicy.CanAccess(project.OrganizationID, authUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return project, nil
}

func (s *ProjectServiceImpl) Create(request *requests.ProjectCreateRequest, authUser models.AuthUser) (models.Project, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationID, authUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.createForbidden")
	}

	err := s.validateNameForDuplication(request.Name, request.OrganizationID)
	if err != nil {
		return models.Project{}, err
	}

	project := models.Project{
		Name:           request.Name,
		OrganizationID: request.OrganizationID,
		DBName:         s.generateDBName(),
		CreatedBy:      authUser.ID,
		UpdatedBy:      authUser.ID,
	}

	_, err = s.projectRepo.Create(&project)
	if err != nil {
		return models.Project{}, err
	}

	err = s.databaseRepo.Create(project.DBName)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (s *ProjectServiceImpl) Update(projectID uuid.UUID, authUser models.AuthUser, request *requests.ProjectUpdateRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationID, authUser) {
		return &models.Project{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = utils.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	project.UpdatedAt = time.Now()
	project.UpdatedBy = authUser.ID

	err = s.validateNameForDuplication(request.Name, project.OrganizationID)
	if err != nil {
		return &models.Project{}, err
	}

	return s.projectRepo.Update(&project)
}

func (s *ProjectServiceImpl) Delete(projectID uuid.UUID, authUser models.AuthUser) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(project.OrganizationID, authUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	err = s.databaseRepo.DropIfExists(project.DBName)
	if err != nil {
		return false, err
	}

	return s.projectRepo.Delete(projectID)
}

func (s *ProjectServiceImpl) generateDBName() string {
	return strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "")
}

func (s *ProjectServiceImpl) validateNameForDuplication(name string, organizationID uuid.UUID) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
