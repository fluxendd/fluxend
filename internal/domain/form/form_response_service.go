package form

import (
	form2 "fluxton/internal/api/dto/form"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FormResponseService interface {
	List(formUUID uuid.UUID, authUser models.AuthUser) ([]FormResponse, error)
	GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser models.AuthUser) (*FormResponse, error)
	Create(formUUID uuid.UUID, request *form2.CreateResponseRequest, authUser models.AuthUser) (FormResponse, error)
	Delete(formUUID, formResponseUUID uuid.UUID, authUser models.AuthUser) error
}

type FormResponseServiceImpl struct {
	projectPolicy              *project.ProjectPolicy
	formFieldValidationService FormFieldValidationService
	formRepo                   *repositories2.FormRepository
	formFieldRepo              *repositories2.FormFieldRepository
	projectRepo                *repositories2.ProjectRepository
	formResponseRepo           *repositories2.FormResponseRepository
}

func NewFormResponseService(injector *do.Injector) (FormResponseService, error) {
	policy := do.MustInvoke[*project.ProjectPolicy](injector)
	formFieldValidationService := do.MustInvoke[FormFieldValidationService](injector)
	formRepo := do.MustInvoke[*repositories2.FormRepository](injector)
	formFieldRepo := do.MustInvoke[*repositories2.FormFieldRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)
	formResponseRepo := do.MustInvoke[*repositories2.FormResponseRepository](injector)

	return &FormResponseServiceImpl{
		projectPolicy:              policy,
		formFieldValidationService: formFieldValidationService,
		formRepo:                   formRepo,
		formFieldRepo:              formFieldRepo,
		projectRepo:                projectRepo,
		formResponseRepo:           formResponseRepo,
	}, nil
}

func (s *FormResponseServiceImpl) List(formUUID uuid.UUID, authUser models.AuthUser) ([]FormResponse, error) {
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

func (s *FormResponseServiceImpl) GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser models.AuthUser) (*FormResponse, error) {
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

func (s *FormResponseServiceImpl) Create(formUUID uuid.UUID, request *form2.CreateResponseRequest, authUser models.AuthUser) (FormResponse, error) {
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

	var formFieldResponses []FormFieldResponse

	for _, formField := range formFields {
		if _, ok := request.Response[formField.Label]; !ok {
			return FormResponse{}, errors.NewUnprocessableError("formResponse.error.missingField")
		}

		currentFieldValue := request.Response[formField.Label].(string)
		validationErr := s.formFieldValidationService.Validate(currentFieldValue, formField)
		if validationErr != nil {
			return FormResponse{}, validationErr
		}

		formFieldResponses = append(formFieldResponses, FormFieldResponse{
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

func (s *FormResponseServiceImpl) Delete(formUUID, formResponseUUID uuid.UUID, authUser models.AuthUser) error {
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
