package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/form_requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FormFieldService interface {
	List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormFiled, error)
	GetByUUID(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.FormFiled, error)
	Create(formUUID, projectUUID uuid.UUID, request *form_requests.CreateFieldRequest, authUser models.AuthUser) (models.FormFiled, error)
	Update(formUUID, projectUUID uuid.UUID, authUser models.AuthUser, request *form_requests.CreateFieldRequest) (*models.FormFiled, error)
	Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FormFieldServiceImpl struct {
	projectPolicy *policies.ProjectPolicy
	formRepo      *repositories.FormRepository
	formFieldRepo *repositories.FormFieldRepository
	projectRepo   *repositories.ProjectRepository
}

func NewFormFieldService(injector *do.Injector) (FormFieldService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	formRepo := do.MustInvoke[*repositories.FormRepository](injector)
	formFieldRepo := do.MustInvoke[*repositories.FormFieldRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)

	return &FormFieldServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		formFieldRepo: formFieldRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FormFieldServiceImpl) List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormFiled, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return []models.FormFiled{}, err
	}

	if !exists {
		return []models.FormFiled{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.FormFiled{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.FormFiled{}, errs.NewForbiddenError("formField.error.listForbidden")
	}

	return s.formFieldRepo.ListForForm(formUUID)
}

func (s *FormFieldServiceImpl) GetByUUID(formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.FormFiled, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.FormFiled{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.FormFiled{}, errs.NewForbiddenError("formField.error.viewForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(formFieldUUID)
	if err != nil {
		return models.FormFiled{}, err
	}

	return formField, nil
}

func (s *FormFieldServiceImpl) Create(formUUID, projectUUID uuid.UUID, request *form_requests.CreateFieldRequest, authUser models.AuthUser) (models.FormFiled, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return models.FormFiled{}, err
	}

	if !exists {
		return models.FormFiled{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.FormFiled{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.FormFiled{}, errs.NewForbiddenError("formField.error.createForbidden")
	}

	err = s.validateLabelForDuplication(request.Label, formUUID)
	if err != nil {
		return models.FormFiled{}, err
	}

	formField := models.FormFiled{
		FormUuid:    formUUID,
		Label:       request.Label,
		Description: request.Description,
		Type:        request.Type,
		IsRequired:  request.IsRequired,
		Options:     request.Options,
	}

	_, err = s.formFieldRepo.Create(&formField)
	if err != nil {
		return models.FormFiled{}, err
	}

	return formField, nil
}

func (s *FormFieldServiceImpl) Update(formUUID, formFieldUUID uuid.UUID, authUser models.AuthUser, request *form_requests.CreateFieldRequest) (*models.FormFiled, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return &models.FormFiled{}, err
	}

	if !exists {
		return &models.FormFiled{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(formUUID)
	if err != nil {
		return &models.FormFiled{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.FormFiled{}, errs.NewForbiddenError("formField.error.updateForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(formFieldUUID)
	if err != nil {
		return &models.FormFiled{}, err
	}

	err = utils.PopulateModel(&formField, request)
	if err != nil {
		return nil, err
	}

	formField.UpdatedAt = time.Now()

	err = s.validateLabelForDuplication(request.Label, formField.FormUuid)
	if err != nil {
		return &models.FormFiled{}, err
	}

	return s.formFieldRepo.Update(&formField)
}

func (s *FormFieldServiceImpl) Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("formField.error.deleteForbidden")
	}

	return s.formRepo.Delete(formUUID)
}

func (s *FormFieldServiceImpl) validateLabelForDuplication(name string, formUUID uuid.UUID) error {
	exists, err := s.formFieldRepo.ExistsByLabelForForm(name, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("formField.error.duplicateName")
	}

	return nil
}
