package services

import (
	"fluxton/models"
	"fluxton/pkg/errors"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/form_requests"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FormResponseService interface {
	List(formUUID uuid.UUID, authUser models.AuthUser) ([]models.FormResponse, error)
	GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser models.AuthUser) (*models.FormResponse, error)
	Create(formUUID uuid.UUID, request *form_requests.CreateResponseRequest, authUser models.AuthUser) (models.FormResponse, error)
	Delete(formUUID, formResponseUUID uuid.UUID, authUser models.AuthUser) error
}

type FormResponseServiceImpl struct {
	projectPolicy              *policies.ProjectPolicy
	formFieldValidationService FormFieldValidationService
	formRepo                   *repositories.FormRepository
	formFieldRepo              *repositories.FormFieldRepository
	projectRepo                *repositories.ProjectRepository
	formResponseRepo           *repositories.FormResponseRepository
}

func NewFormResponseService(injector *do.Injector) (FormResponseService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	formFieldValidationService := do.MustInvoke[FormFieldValidationService](injector)
	formRepo := do.MustInvoke[*repositories.FormRepository](injector)
	formFieldRepo := do.MustInvoke[*repositories.FormFieldRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	formResponseRepo := do.MustInvoke[*repositories.FormResponseRepository](injector)

	return &FormResponseServiceImpl{
		projectPolicy:              policy,
		formFieldValidationService: formFieldValidationService,
		formRepo:                   formRepo,
		formFieldRepo:              formFieldRepo,
		projectRepo:                projectRepo,
		formResponseRepo:           formResponseRepo,
	}, nil
}

func (s *FormResponseServiceImpl) List(formUUID uuid.UUID, authUser models.AuthUser) ([]models.FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return []models.FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.FormResponse{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.FormResponse{}, errors.NewForbiddenError("formFieldResponse.error.listForbidden")
	}

	return s.formResponseRepo.ListForForm(formUUID)
}

func (s *FormResponseServiceImpl) GetByUUID(formResponseUUID, formUUID uuid.UUID, authUser models.AuthUser) (*models.FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return &models.FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return &models.FormResponse{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return &models.FormResponse{}, errors.NewForbiddenError("formFieldResponse.error.showForbidden")
	}

	return s.formResponseRepo.GetByUUID(formResponseUUID)
}

func (s *FormResponseServiceImpl) Create(formUUID uuid.UUID, request *form_requests.CreateResponseRequest, authUser models.AuthUser) (models.FormResponse, error) {
	projectUUID, err := s.formRepo.GetProjectUUIDByFormUUID(formUUID)
	if err != nil {
		return models.FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.FormResponse{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.FormResponse{}, errors.NewForbiddenError("formResponse.error.createForbidden")
	}

	formFields, err := s.formFieldRepo.ListForForm(formUUID)
	if err != nil {
		return models.FormResponse{}, err
	}

	formResponse := models.FormResponse{
		FormUuid: formUUID,
	}

	var formFieldResponses []models.FormFieldResponse

	for _, formField := range formFields {
		if _, ok := request.Response[formField.Label]; !ok {
			return models.FormResponse{}, errors.NewUnprocessableError("formResponse.error.missingField")
		}

		currentFieldValue := request.Response[formField.Label].(string)
		validationErr := s.formFieldValidationService.Validate(currentFieldValue, formField)
		if validationErr != nil {
			return models.FormResponse{}, validationErr
		}

		formFieldResponses = append(formFieldResponses, models.FormFieldResponse{
			FormFieldUuid: formField.Uuid,
			Value:         currentFieldValue,
		})
	}

	_, err = s.formResponseRepo.Create(&formResponse, &formFieldResponses)
	if err != nil {
		return models.FormResponse{}, err
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
