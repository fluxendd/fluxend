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
)

type ProjectService interface {
	List(paginationParams utils.PaginationParams, organizationId, authenticatedUserId uuid.UUID) ([]models.Project, error)
	GetByID(projectId, organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Project, error)
	Create(request *requests.ProjectCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Project, error)
	Update(projectId uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.ProjectCreateRequest) (*models.Project, error)
	Delete(projectId uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error)
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

func (s *ProjectServiceImpl) List(paginationParams utils.PaginationParams, organizationId, authenticatedUserId uuid.UUID) ([]models.Project, error) {
	if !s.projectPolicy.CanList(organizationId, authenticatedUserId) {
		return []models.Project{}, errs.NewForbiddenError("project.error.listForbidden")
	}

	return s.projectRepo.ListForUser(paginationParams, authenticatedUserId)
}

func (s *ProjectServiceImpl) GetByID(projectId, organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Project, error) {
	if !s.projectPolicy.CanView(organizationId, authenticatedUser) {
		return models.Project{}, errs.NewForbiddenError("project.error.viewForbidden")
	}

	return s.projectRepo.GetByID(projectId)
}

func (s *ProjectServiceImpl) Create(request *requests.ProjectCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Project, error) {
	if !s.projectPolicy.CanCreate(request.OrganizationID, authenticatedUser) {
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

func (s *ProjectServiceImpl) Update(projectId uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.ProjectCreateRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(projectId)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanUpdate(projectId, authenticatedUser) {
		return &models.Project{}, errs.NewForbiddenError("project.error.updateForbidden")
	}

	project.OrganizationID = request.OrganizationID
	err = utils.PopulateModel(&project, request)
	if err != nil {
		return nil, err
	}

	err = s.validateNameForDuplication(request.Name, request.OrganizationID)
	if err != nil {
		return &models.Project{}, err
	}

	return s.projectRepo.Update(projectId, &project)
}

func (s *ProjectServiceImpl) Delete(projectId uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error) {
	if !s.projectPolicy.CanUpdate(projectId, authenticatedUser) {
		return false, errs.NewForbiddenError("project.error.updateForbidden")
	}

	return s.projectRepo.Delete(projectId)
}

func (s *ProjectServiceImpl) generateDBName() string {
	return strings.ReplaceAll(strings.ToLower(uuid.New().String()), "-", "")
}

func (s *ProjectServiceImpl) validateNameForDuplication(name string, organizationId uuid.UUID) error {
	exists, err := s.projectRepo.ExistsByNameForOrganization(name, organizationId)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("project.error.duplicateName")
	}

	return nil
}
