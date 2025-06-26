package form

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/project"
	"fluxend/internal/domain/shared"
	"fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(paginationParams shared.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Form, error)
	GetByUUID(formUUID uuid.UUID, authUser auth.User) (Form, error)
	Create(request *CreateFormInput, authUser auth.User) (Form, error)
	Update(formUUID uuid.UUID, authUser auth.User, request *CreateFormInput) (*Form, error)
	Delete(formUUID uuid.UUID, authUser auth.User) (bool, error)
}

type ServiceImpl struct {
	projectPolicy *project.Policy
	formRepo      Repository
	projectRepo   project.Repository
}

func NewFormService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formRepo := do.MustInvoke[Repository](injector)
	projectRepo := do.MustInvoke[project.Repository](injector)

	return &ServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *ServiceImpl) List(paginationParams shared.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Form, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("form.error.listForbidden")
	}

	return s.formRepo.ListForProject(paginationParams, projectUUID)
}

func (s *ServiceImpl) GetByUUID(formUUID uuid.UUID, authUser auth.User) (Form, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedForm.ProjectUuid)
	if err != nil {
		return Form{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Form{}, errors.NewForbiddenError("form.error.viewForbidden")
	}

	return fetchedForm, nil
}

func (s *ServiceImpl) Create(request *CreateFormInput, authUser auth.User) (Form, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return Form{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return Form{}, errors.NewForbiddenError("form.error.createForbidden")
	}

	if err = s.validateNameForDuplication(request.Name, request.ProjectUUID); err != nil {
		return Form{}, err
	}

	formInput := Form{
		ProjectUuid: request.ProjectUUID,
		Name:        request.Name,
		Description: request.Description,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
		CreatedAt:   time.Now(),
	}

	_, err = s.formRepo.Create(&formInput)
	if err != nil {
		return Form{}, err
	}

	return formInput, nil
}

func (s *ServiceImpl) Update(formUUID uuid.UUID, authUser auth.User, request *CreateFormInput) (*Form, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedForm.ProjectUuid)
	if err != nil {
		return &Form{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &Form{}, errors.NewForbiddenError("form.error.updateForbidden")
	}

	if err = fetchedForm.PopulateModel(&fetchedForm, request); err != nil {
		return nil, err
	}

	fetchedForm.UpdatedAt = time.Now()
	fetchedForm.UpdatedBy = authUser.Uuid

	if err = s.validateNameForDuplication(request.Name, fetchedForm.ProjectUuid); err != nil {
		return &Form{}, err
	}

	return s.formRepo.Update(&fetchedForm)
}

func (s *ServiceImpl) Delete(formUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(fetchedForm.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("form.error.deleteForbidden")
	}

	return s.formRepo.Delete(formUUID)
}

func (s *ServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.formRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("form.error.duplicateName")
	}

	return nil
}
