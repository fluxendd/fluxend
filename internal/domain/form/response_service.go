package form

import (
	"fluxton/internal/api/dto/form"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ResponseService interface {
	List(formUUID uuid.UUID, authUser auth.User) ([]FormResponse, error)
	GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser auth.User) (*FormResponse, error)
	Create(formUUID uuid.UUID, request *form.CreateResponseRequest, authUser auth.User) (FormResponse, error)
	Delete(formUUID, formResponseUUID uuid.UUID, authUser auth.User) error
}

type ResponseServiceImpl struct {
	projectPolicy              *project.Policy
	formFieldValidationService FieldValidationService
	formRepo                   *Repository
	formFieldRepo              *FieldRepository
	projectRepo                *project.Repository
	formResponseRepo           *FieldResponseRepository
}

func NewFormResponseService(injector *do.Injector) (ResponseService, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formFieldValidationService := do.MustInvoke[FieldValidationService](injector)
	formRepo := do.MustInvoke[*Repository](injector)
	formFieldRepo := do.MustInvoke[*FieldRepository](injector)
	projectRepo := do.MustInvoke[*project.Repository](injector)
	formResponseRepo := do.MustInvoke[*FieldResponseRepository](injector)

	return &ResponseServiceImpl{
		projectPolicy:              policy,
		formFieldValidationService: formFieldValidationService,
		formRepo:                   formRepo,
		formFieldRepo:              formFieldRepo,
		projectRepo:                projectRepo,
		formResponseRepo:           formResponseRepo,
	}, nil
}

func (s *ResponseServiceImpl) List(formUUID uuid.UUID, authUser auth.User) ([]FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []FormResponse{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []FormResponse{}, errors.NewForbiddenError("formFieldResponse.error.listForbidden")
	}

	return s.formResponseRepo.ListForForm(formUUID)
}

func (s *ResponseServiceImpl) GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser auth.User) (*FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return &FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return &FormResponse{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return &FormResponse{}, errors.NewForbiddenError("formFieldResponse.error.showForbidden")
	}

	return s.formResponseRepo.GetByUUID(formResponseUUID)
}

func (s *ResponseServiceImpl) Create(formUUID uuid.UUID, request *form.CreateResponseRequest, authUser auth.User) (FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return FormResponse{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return FormResponse{}, errors.NewForbiddenError("formResponse.error.createForbidden")
	}

	formFields, err := s.formFieldRepo.ListForForm(formUUID)
	if err != nil {
		return FormResponse{}, err
	}

	formResponse := FormResponse{
		FormUuid: formUUID,
	}

	var formFieldResponses []FieldResponse

	for _, formField := range formFields {
		if _, ok := request.Response[formField.Label]; !ok {
			return FormResponse{}, errors.NewUnprocessableError("formResponse.error.missingField")
		}

		currentFieldValue := request.Response[formField.Label].(string)
		validationErr := s.formFieldValidationService.Validate(currentFieldValue, formField)
		if validationErr != nil {
			return FormResponse{}, validationErr
		}

		formFieldResponses = append(formFieldResponses, FieldResponse{
			FormFieldUuid: formField.Uuid,
			Value:         currentFieldValue,
		})
	}

	_, err = s.formResponseRepo.Create(&formResponse, &formFieldResponses)
	if err != nil {
		return FormResponse{}, err
	}

	return formResponse, nil
}

func (s *ResponseServiceImpl) Delete(formUUID, formResponseUUID uuid.UUID, authUser auth.User) error {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return errors.NewForbiddenError("form.error.deleteForbidden")
	}

	return s.formResponseRepo.Delete(formResponseUUID)
}
