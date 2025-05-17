package form

import (
	"fluxton/internal/api/dto"
	form2 "fluxton/internal/api/dto/form"
	repositories2 "fluxton/internal/database/repositories"
	"fluxton/internal/domain/project"
	"fluxton/models"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type FormService interface {
	List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]Form, error)
	GetByUUID(formUUID uuid.UUID, authUser models.AuthUser) (Form, error)
	Create(request *form2.CreateRequest, authUser models.AuthUser) (Form, error)
	Update(formUUID uuid.UUID, authUser models.AuthUser, request *form2.CreateRequest) (*Form, error)
	Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error)
}

type FormServiceImpl struct {
	projectPolicy *project.Policy
	formRepo      *repositories2.FormRepository
	projectRepo   *repositories2.ProjectRepository
}

func NewFormService(injector *do.Injector) (FormService, error) {
	policy := do.MustInvoke[*project.Policy](injector)
	formRepo := do.MustInvoke[*repositories2.FormRepository](injector)
	projectRepo := do.MustInvoke[*repositories2.ProjectRepository](injector)

	return &FormServiceImpl{
		projectPolicy: policy,
		formRepo:      formRepo,
		projectRepo:   projectRepo,
	}, nil
}

func (s *FormServiceImpl) List(paginationParams dto.PaginationParams, projectUUID uuid.UUID, authUser models.AuthUser) ([]Form, error) {
	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(projectUUID)
	if err != nil {
		return nil, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("form.error.listForbidden")
	}

	return s.formRepo.ListForProject(paginationParams, projectUUID)
}

func (s *FormServiceImpl) GetByUUID(formUUID uuid.UUID, authUser models.AuthUser) (Form, error) {
	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return form.Form{}, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return form.Form{}, err
	}

	if !s.projectPolicy.CanAccess(organizationUUID, authUser) {
		return form.Form{}, errors.NewForbiddenError("form.error.viewForbidden")
	}

	return form, nil
}

func (s *FormServiceImpl) Create(request *form2.CreateRequest, authUser models.AuthUser) (Form, error) {
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
		return form.Form{}, err
	}

	return form, nil
}

func (s *FormServiceImpl) Update(formUUID uuid.UUID, authUser models.AuthUser, request *form2.CreateRequest) (*Form, error) {
	form, err := s.formRepo.GetByUUID(formUUID)
	if err != nil {
		return nil, err
	}

	organizationUUID, err := s.projectRepo.GetOrganizationUUIDByProjectUUID(form.ProjectUuid)
	if err != nil {
		return &form.Form{}, err
	}

	if !s.projectPolicy.CanUpdate(organizationUUID, authUser) {
		return &form.Form{}, errors.NewForbiddenError("form.error.updateForbidden")
	}

	err = form.PopulateModel(&form, request)
	if err != nil {
		return nil, err
	}

	form.UpdatedAt = time.Now()
	form.UpdatedBy = authUser.Uuid

	err = s.validateNameForDuplication(request.Name, form.ProjectUuid)
	if err != nil {
		return &form.Form{}, err
	}

	return s.formRepo.Update(&form)
}

func (s *FormServiceImpl) Delete(formUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	form, err := s.formRepo.GetByUUID(formUUID)
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
