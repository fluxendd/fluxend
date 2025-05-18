package form

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/form"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/project"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Form, error)
	GetByUUID(formUUID uuid.UUID, authUser auth.User) (Form, error)
	Create(request *form.CreateRequest, authUser auth.User) (Form, error)
	Update(formUUID uuid.UUID, authUser auth.User, request *form.CreateRequest) (*Form, error)
	Delete(formUUID uuid.UUID, authUser auth.User) (bool, error)
}

type FormServiceImpl struct {
	projectPolicy *project.Policy
	formRepo      *Repository
	projectRepo   *project.Repository
}

func NewFormService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formRepo := do.MustInvoke[*Repository](injector)
	projectRepo := do.MustInvoke[*project.Repository](injector)

	return &FormServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FormServiceImpl) List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser auth.User) ([]Form, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("form.error.listForbidden")
	}

	return s.formRepo.ListForProject(paginationParams, projectUUID)
}

func (s *FormServiceImpl) GetByUUID(formUUID uuid.UUID, authUser auth.User) (Form, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return Form{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return Form{}, errors.NewForbiddenError("form.error.viewForbidden")
	}

	return fetchedForm, nil
}

func (s *FormServiceImpl) Create(request *form.CreateRequest, authUser auth.User) (Form, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(request.ProjectUUID)
	if err != nil {
		return Form{}, err
	}

	if !s.projectPolicy.CanCreate(organizationUUID, authUser) {
		return Form{}, errors.NewForbiddenError("form.error.createForbidden")
	}

	err = s.validateNameForDuplication(request.Name, request.ProjectUUID)
	if err != nil {
		return Form{}, err
	}

	form := Form{
		ProjectUuid: request.ProjectUUID,
		Name:        request.Name,
		Description: request.Description,
		CreatedBy:   authUser.Uuid,
		UpdatedBy:   authUser.Uuid,
		CreatedAt:   time.Now(),
	}

	_, err = s.formRepo.Create(&form)
	if err != nil {
		return Form{}, err
	}

	return form, nil
}

func (s *FormServiceImpl) Update(formUUID uuid.UUID, authUser auth.User, request *form.CreateRequest) (*Form, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return &Form{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &Form{}, errors.NewForbiddenError("form.error.updateForbidden")
	}

	/*err = form.PopulateModel(&fetchedForm, request)
	if err != nil {
		return nil, err
	}*/

	fetchedForm.UpdatedAt = time.Now()
	fetchedForm.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, fetchedForm.ProjectUuid)
	if err != nil {
		return &Form{}, err
	}

	return s.formRepo.Update(&fetchedForm)
}

func (s *FormServiceImpl) Delete(formUUID uuid.UUID, authUser auth.User) (bool, error) {
	fetchedForm, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return false, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return false, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("form.error.deleteForbidden")
	}

	return s.formRepo.Delete(formUUID)
}

func (s *FormServiceImpl) validateNameForDuplication(name string, projectUUID uuid.UUID) error {
	exists, err := s.formRepo.ExistsByNameForProject(name, projectUUID)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewUnprocessableError("form.error.duplicateName")
	}

	return nil
}
