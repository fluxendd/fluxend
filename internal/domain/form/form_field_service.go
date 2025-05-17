package form

import (
	form2 "fluxton/internal/api/dto/form"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FormFieldService interface {
	List(formUUID uuid.UUID, authUser models.AuthUser) ([]FormField, error)
	GetByUUID(formUUID uuid.UUID, authUser models.AuthUser) (FormField, error)
	CreateMany(formUUID uuid.UUID, request *form2.CreateFormFieldsRequest, authUser models.AuthUser) ([]FormField, error)
	Update(formUUID, fieldUUID uuid.UUID, authUser models.AuthUser, request *form2.UpdateFormFieldRequest) (*FormField, error)
	Delete(formUUID, fieldUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FormFieldServiceImpl struct {
	projectPolicy *project.Policy
	formRepo      *repositories2.FormRepository
	formFieldRepo *repositories2.FormFieldRepository
	projectRepo   *repositories2.ProjectRepository
}

func NewFormFieldService(injector *do.Injector) (FormFieldService, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formRepo := do.MustInvoke[*repositories2.FormRepository](injector)
	formFieldRepo := do.MustInvoke[*repositories2.FormFieldRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &FormFieldServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		formFieldRepo: formFieldRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FormFieldServiceImpl) List(formUUID uuid.UUID, authUser models.AuthUser) ([]FormField, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []FormField{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []FormField{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []FormField{}, errors.NewForbiddenError("formField.error.listForbidden")
	}

	return s.formFieldRepo.ListForForm(formUUID)
}

func (s *FormFieldServiceImpl) GetByUUID(fieldUUID uuid.UUID, authUser models.AuthUser) (FormField, error) {
	formField, err := s.formFieldRepo.GetByUUID(fieldUUID)
	if err != nil {
		return FormField{}, err
	}

	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formField.FormUuid)
	if err != nil {
		return FormField{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return FormField{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return FormField{}, errors.NewForbiddenError("formField.error.viewForbidden")
	}

	return formField, nil
}

func (s *FormFieldServiceImpl) CreateMany(formUUID uuid.UUID, request *form2.CreateFormFieldsRequest, authUser models.AuthUser) ([]FormField, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []FormField{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []FormField{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return []FormField{}, errors.NewForbiddenError("formField.error.createForbidden")
	}

	err = s.validateManyForLabelDuplication(request, formUUID)
	if err != nil {
		return []FormField{}, err
	}

	formFields := make([]FormField, len(request.Fields))
	for i, field := range request.Fields {
		currentField := FormField{}
		err := currentField.PopulateModel(&currentField, field)
		if err != nil {
			return []FormField{}, err
		}

		formFields[i] = currentField
	}

	createdFields, err := s.formFieldRepo.CreateMany(formFields, formUUID)
	if err != nil {
		return []FormField{}, err
	}

	return createdFields, nil
}

func (s *FormFieldServiceImpl) Update(formUUID, fieldUUID uuid.UUID, authUser models.AuthUser, request *form2.UpdateFormFieldRequest) (*FormField, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return &FormField{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return &FormField{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &FormField{}, errors.NewForbiddenError("formField.error.updateForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(fieldUUID)
	if err != nil {
		return &FormField{}, err
	}

	err = formField.PopulateModel(&formField, request.FieldRequest)
	if err != nil {
		return nil, err
	}

	formField.UpdatedAt = time.Now()

	err = s.validateOneForLabelDuplication(request.Label, formField.FormUuid)
	if err != nil {
		return &FormField{}, err
	}

	return s.formFieldRepo.Update(&formField)
}

func (s *FormFieldServiceImpl) Delete(formUUID, fieldUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("formField.error.deleteForbidden")
	}

	return s.formFieldRepo.Delete(fieldUUID)
}

func (s *FormFieldServiceImpl) validateManyForLabelDuplication(request *form2.CreateFormFieldsRequest, formUUID uuid.UUID) error {
	labels := make([]string, len(request.Fields))

	for i, field := range request.Fields {
		labels[i] = field.Label
	}

	exists, err := s.formFieldRepo.ExistsByAnyLabelForForm(labels, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("formField.error.someDuplicateLabels")
	}

	return nil
}

func (s *FormFieldServiceImpl) validateOneForLabelDuplication(label string, formUUID uuid.UUID) error {
	exists, err := s.formFieldRepo.ExistsByLabelForForm(label, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("formField.error.duplicateLabel")
	}

	return nil
}
