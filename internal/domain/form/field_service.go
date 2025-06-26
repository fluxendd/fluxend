package form

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FieldService interface {
	List(formUUID uuid.UUID, authUser auth.User) ([]Field, error)
	GetByUUID(formUUID uuid.UUID, authUser auth.User) (Field, error)
	CreateMany(formUUID uuid.UUID, request *CreateFormFieldsInput, authUser auth.User) ([]Field, error)
	Update(formUUID, fieldUUID uuid.UUID, authUser auth.User, request *UpdateFormFieldsInput) (*Field, error)
	Delete(formUUID, fieldUUID uuid.UUID, authUser auth.User) (bool, error)
}

type FieldServiceImpl struct {
	projectPolicy *project.Policy
	formRepo      Repository
	formFieldRepo FieldRepository
	projectRepo   project.Repository
}

func NewFieldService(injector *do.Injector) (FieldService, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formRepo := do.MustInvoke[Repository](injector)
	formFieldRepo := do.MustInvoke[FieldRepository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &FieldServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		formFieldRepo: formFieldRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FieldServiceImpl) List(formUUID uuid.UUID, authUser auth.User) ([]Field, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []Field{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Field{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []Field{}, errors.NewForbiddenError("formField.error.listForbidden")
	}

	return s.formFieldRepo.ListForForm(formUUID)
}

func (s *FieldServiceImpl) GetByUUID(fieldUUID uuid.UUID, authUser auth.User) (Field, error) {
	formField, err := s.formFieldRepo.GetByUUID(fieldUUID)
	if err != nil {
		return Field{}, err
	}

	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formField.FormUuid)
	if err != nil {
		return Field{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return Field{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Field{}, errors.NewForbiddenError("formField.error.viewForbidden")
	}

	return formField, nil
}

func (s *FieldServiceImpl) CreateMany(formUUID uuid.UUID, request *CreateFormFieldsInput, authUser auth.User) ([]Field, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []Field{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []Field{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return []Field{}, errors.NewForbiddenError("formField.error.createForbidden")
	}

	if err = s.validateManyForLabelDuplication(request, formUUID); err != nil {
		return []Field{}, err
	}

	formFields := make([]Field, len(request.Fields))
	for i, field := range request.Fields {
		currentField := Field{}
		err := currentField.PopulateModel(&currentField, field)
		if err != nil {
			return []Field{}, err
		}

		formFields[i] = currentField
	}

	createdFields, err := s.formFieldRepo.CreateMany(formFields, formUUID)
	if err != nil {
		return []Field{}, err
	}

	return createdFields, nil
}

func (s *FieldServiceImpl) Update(formUUID, fieldUUID uuid.UUID, authUser auth.User, request *UpdateFormFieldsInput) (*Field, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return &Field{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return &Field{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &Field{}, errors.NewForbiddenError("formField.error.updateForbidden")
	}

	formField, err := s.formFieldRepo.GetByUUID(fieldUUID)
	if err != nil {
		return &Field{}, err
	}

	if err = formField.PopulateModel(&formField, request.FieldInput); err != nil {
		return nil, err
	}

	formField.UpdatedAt = time.Now()

	if err = s.validateOneForLabelDuplication(request.Label, formField.FormUuid); err != nil {
		return &Field{}, err
	}

	return s.formFieldRepo.Update(&formField)
}

func (s *FieldServiceImpl) Delete(formUUID, fieldUUID uuid.UUID, authUser auth.User) (bool, error) {
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

func (s *FieldServiceImpl) validateManyForLabelDuplication(request *CreateFormFieldsInput, formUUID uuid.UUID) error {
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

func (s *FieldServiceImpl) validateOneForLabelDuplication(label string, formUUID uuid.UUID) error {
	exists, err := s.formFieldRepo.ExistsByLabelForForm(label, formUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("formField.error.duplicateLabel")
	}

	return nil
}
