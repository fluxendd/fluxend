package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/requests/form_requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FormService interface {
	List(paginationParams requests.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Form, error)
	GetByUUID(formUUID uuid.UUID, authUser models.AuthUser) (models.Form, error)
	Create(request *form_requests.CreateRequest, authUser models.AuthUser) (models.Form, error)
	Update(formUUID uuid.UUID, authUser models.AuthUser, request *form_requests.CreateRequest) (*models.Form, error)
	Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FormServiceImpl struct {
	projectPolicy  *policies.ProjectPolicy
	settingService SettingService
	formRepo       *repositories.FormRepository
	projectRepo    *repositories.ProjectRepository
}

func NewFormService(injector *do.Injector) (FormService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	settingService := do.MustInvoke[SettingService](injector)
	formRepo := do.MustInvoke[*repositories.FormRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FormServiceImpl{
		projectPolicy:  policy,
		settingService: settingService,
		formRepo:       formRepo,
		projectRepo:    projectRepo,
	}, nil
}

func (s *FormServiceImpl) List(paginationParams requests.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.Form, error) {
	err := s.validateFormsEnabled()
	if err != nil {
		return []models.Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.Form{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.Form{}, errs.NewForbiddenError("form.error.listForbidden")
	}

	return s.formRepo.ListForProject(paginationParams, projectUUID)
}

func (s *FormServiceImpl) GetByUUID(formUUID uuid.UUID, authUser models.AuthUser) (models.Form, error) {
	err := s.validateFormsEnabled()
	if err != nil {
		return models.Form{}, err
	}

	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return models.Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return models.Form{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.Form{}, errs.NewForbiddenError("form.error.viewForbidden")
	}

	return form, nil
}

func (s *FormServiceImpl) Create(request *form_requests.CreateRequest, authUser models.AuthUser) (models.Form, error) {
	err := s.validateFormsEnabled()
	if err != nil {
		return models.Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return models.Form{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.Form{}, errs.NewForbiddenError("form.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, request.ProjectUUID)
	if err != nil {
		return models.Form{}, err
	}

	form := models.Form{
		ProjectUuid: request.ProjectUUID,
		Name:        request.Name,
		Description: request.Description,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
		CreatedAt:   time.Now(),
	}

	_, err = s.formRepo.Create(&form)
	if err != nil {
		return models.Form{}, err
	}

	return form, nil
}

func (s *FormServiceImpl) Update(formUUID uuid.UUID, authUser models.AuthUser, request *form_requests.CreateRequest) (*models.Form, error) {
	err := s.validateFormsEnabled()
	if err != nil {
		return &models.Form{}, err
	}

	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return &models.Form{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.Form{}, errs.NewForbiddenError("form.error.updateForbidden")
	}

	err = utils.PopulateModel(&form, request)
	if err != nil {
		return nil, err
	}

	form.UpdatedAt = time.Now()
	form.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, form.ProjectUuid)
	if err != nil {
		return &models.Form{}, err
	}

	return s.formRepo.Update(&form)
}

func (s *FormServiceImpl) Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	err := s.validateFormsEnabled()
	if err != nil {
		return false, err
	}

	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("form.error.deleteForbidden")
	}

	return s.formRepo.Delete(formUUID)
}

func (s *FormServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.formRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("form.error.duplicateName")
	}

	return nil
}

func (s *FormServiceImpl) validateFormsEnabled() error {
	enabled := s.settingService.GetBool("enableForms")
	if !enabled {
		return errs.NewForbiddenError("form.error.formsDisabled")
	}

	return nil
}
