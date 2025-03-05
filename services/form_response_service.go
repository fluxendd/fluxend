package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/form_requests"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type FormResponseService interface {
	List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormResponse, error)
	Create(formUUID, projectUUID uuid.UUID, request *form_requests.CreateResponseRequest, authUser models.AuthUser) (models.FormResponse, error)
}

type FormResponseServiceImpl struct {
	projectPolicy    *policies.ProjectPolicy
	formRepo         *repositories.FormRepository
	formFieldRepo    *repositories.FormFieldRepository
	projectRepo      *repositories.ProjectRepository
	formResponseRepo *repositories.FormResponseRepository
}

func NewFormResponseService(injector *do.Injector) (FormResponseService, error) {
	policy := do.MustInvoke[*policies.ProjectPolicy](injector)
	formRepo := do.MustInvoke[*repositories.FormRepository](injector)
	formFieldRepo := do.MustInvoke[*repositories.FormFieldRepository](injector)
	projectRepo := do.MustInvoke[*repositories.ProjectRepository](injector)
	formResponseRepo := do.MustInvoke[*repositories.FormResponseRepository](injector)

	return &FormResponseServiceImpl{
		projectPolicy:    policy,
		formRepo:         formRepo,
		formFieldRepo:    formFieldRepo,
		projectRepo:      projectRepo,
		formResponseRepo: formResponseRepo,
	}, nil
}

func (s *FormResponseServiceImpl) List(formUUID, projectUUID uuid.UUID, authUser models.AuthUser) ([]models.FormResponse, error) {
	err := s.validateFormExists(formUUID)
	if err != nil {
		return []models.FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return []models.FormResponse{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return []models.FormResponse{}, errs.NewForbiddenError("formFieldResponse.error.listForbidden")
	}

	return s.formResponseRepo.ListForForm(formUUID)
}

func (s *FormResponseServiceImpl) Create(formUUID, projectUUID uuid.UUID, request *form_requests.CreateResponseRequest, authUser models.AuthUser) (models.FormResponse, error) {
	err := s.validateFormExists(formUUID)
	if err != nil {
		return models.FormResponse{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return models.FormResponse{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return models.FormResponse{}, errs.NewForbiddenError("formResponse.error.createForbidden")
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
			return models.FormResponse{}, errs.NewUnprocessableError("formResponse.error.missingField")
		}

		formFieldResponses = append(formFieldResponses, models.FormFieldResponse{
			FormFieldUuid: formField.Uuid,
			Value:         request.Response[formField.Label].(string),
		})
	}

	_, err = s.formResponseRepo.Create(&formResponse, &formFieldResponses)
	if err != nil {
		return models.FormResponse{}, err
	}

	return formResponse, nil
}

func (s *FormResponseServiceImpl) validateFormExists(formUUID uuid.UUID) error {
	formExists, err := s.formRepo.ExistsByUUID(formUUID)
	if err != nil {
		return err
	}

	if !formExists {
		return errs.NewNotFoundError("form.error.notFound")
	}

	return nil
}
