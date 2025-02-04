package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type OrganizationService interface {
	List(paginationParams utils.PaginationParams, authenticatedUserId uuid.UUID) ([]models.Organization, error)
	GetByID(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Organization, error)
	Create(request *requests.OrganizationCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Organization, error)
	Update(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.OrganizationCreateRequest) (*models.Organization, error)
	Delete(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error)
}

type OrganizationServiceImpl struct {
	organizationPolicy *policies.OrganizationPolicy
	organizationRepo   *repositories.OrganizationRepository
}

func NewOrganizationService(injector *do.Injector) (OrganizationService, error) {
	policy := do.MustInvoke[*policies.OrganizationPolicy](injector)
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

	return &OrganizationServiceImpl{
		organizationPolicy: policy,
		organizationRepo:   repo,
	}, nil
}

func (s *OrganizationServiceImpl) List(paginationParams utils.PaginationParams, authenticatedUserId uuid.UUID) ([]models.Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authenticatedUserId)
}

func (s *OrganizationServiceImpl) GetByID(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (models.Organization, error) {
	organization, err := s.organizationRepo.GetByIDForUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return models.Organization{}, err
	}

	if !s.organizationPolicy.CanView(organizationId, authenticatedUser) {
		return models.Organization{}, errs.NewForbiddenError("organization.error.viewForbidden")
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) Create(request *requests.OrganizationCreateRequest, authenticatedUser models.AuthenticatedUser) (models.Organization, error) {
	if !s.organizationPolicy.CanCreate(authenticatedUser) {
		return models.Organization{}, errs.NewForbiddenError("organization.error.createForbidden")
	}

	organization := models.Organization{
		Name: request.Name,
	}

	_, err := s.organizationRepo.Create(&organization, authenticatedUser.ID)
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) Update(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser, request *requests.OrganizationCreateRequest) (*models.Organization, error) {
	organization, err := s.organizationRepo.GetByIDForUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organizationId, authenticatedUser) {
		return &models.Organization{}, errs.NewForbiddenError("organization.error.updateForbidden")
	}

	err = utils.PopulateModel(&organization, request)
	if err != nil {
		return nil, err
	}

	return s.organizationRepo.Update(organizationId, &organization)
}

func (s *OrganizationServiceImpl) Delete(organizationId uuid.UUID, authenticatedUser models.AuthenticatedUser) (bool, error) {
	_, err := s.organizationRepo.GetByIDForUser(organizationId, authenticatedUser.ID)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organizationId, authenticatedUser) {
		return false, errs.NewForbiddenError("organization.error.updateForbidden")
	}

	return s.organizationRepo.Delete(organizationId)
}
