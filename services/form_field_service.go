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
	List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormField, error)
	GetByUUID(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.FormField, error)
	CreateMany(formUUID, projectUUID uuid.UUID, request *form_requests.CreateFormFieldsRequest, authUser models.AuthUser) ([]models.FormField, error)
	Update(formUUID, formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser, request *form_requests.UpdateFormFieldRequest) (*models.FormField, error)
	Delete(formUUID, formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error)
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

func (s *FormFieldServiceImpl) List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormField, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	if !exists {
		return []models.FormField{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.FormField{}, errs.NewForbiddenError("formField.error.listForbidden")
	}

	return s.formFieldRepo.ListForForm(formUUID)
}

func (s *FormFieldServiceImpl) GetByUUID(formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser) (models.FormField, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.FormField{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return models.FormField{}, errs.NewForbiddenError("formField.error.viewForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(formFieldUUID)
	if err != nil {
		return models.FormField{}, err
	}

	return formField, nil
}

func (s *FormFieldServiceImpl) CreateMany(formUUID, projectUUID uuid.UUID, request *form_requests.CreateFormFieldsRequest, authUser models.AuthUser) ([]models.FormField, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	if !exists {
		return []models.FormField{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return []models.FormField{}, errs.NewForbiddenError("formField.error.createForbidden")
	}

	err = s.validateManyForLabelDuplication(request, formUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	formFields := make([]models.FormField, len(request.Fields))
	for i, field := range request.Fields {
		formFields[i] = models.FormField{
			FormUuid:   formUUID,
			Label:      field.Label,
			Type:       field.Type,
			IsRequired: field.IsRequired,
			Options:    field.Options,
		}
	}

	_, err = s.formFieldRepo.CreateMany(formFields, formUUID)
	if err != nil {
		return []models.FormField{}, err
	}

	return formFields, nil
}

func (s *FormFieldServiceImpl) Update(formUUID, formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser, request *form_requests.UpdateFormFieldRequest) (*models.FormField, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return &models.FormField{}, err
	}

	if !exists {
		return &models.FormField{}, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return &models.FormField{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.FormField{}, errs.NewForbiddenError("formField.error.updateForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(formFieldUUID)
	if err != nil {
		return &models.FormField{}, err
	}

	err = utils.PopulateModel(&formField, request)
	if err != nil {
		return nil, err
	}

	formField.UpdatedAt = time.Now()

	err = s.validateOneForLabelDuplication(request.Label, formField.FormUuid)
	if err != nil {
		return &models.FormField{}, err
	}

	return s.formFieldRepo.Update(&formField)
}

func (s *FormFieldServiceImpl) Delete(formUUID, formFieldUUID, projectUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	exists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errs.NewNotFoundError("form.error.notFound")
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("formField.error.deleteForbidden")
	}

	return s.formFieldRepo.Delete(formFieldUUID)
}

func (s *FormFieldServiceImpl) validateManyForLabelDuplication(request *form_requests.CreateFormFieldsRequest, formUUID uuid.UUID) error {
	labels := make([]string, len(request.Fields))

	for i, field := range request.Fields {
		labels[i] = field.Label
	}

	exists, err := s.formFieldRepo.ExistsByAnyLabelForForm(labels, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("formField.error.someDuplicateLabels")
	}

	return nil
}

func (s *FormFieldServiceImpl) validateOneForLabelDuplication(label string, formUUID uuid.UUID) error {
	exists, err := s.formFieldRepo.ExistsByLabelForForm(label, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errs.NewUnprocessableError("formField.error.duplicateLabel")
	}

	return nil
}
