package services

import (
	"github.com/samber/do"
	"myapp/errs"
	"myapp/models"
	"myapp/policies"
	"myapp/repositories"
	"myapp/requests"
	"myapp/utils"
)

type OrganizationService interface {
	List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Organization, error)
	GetByID(id, authenticatedUserId uint) (models.Organization, error)
	Create(request *requests.OrganizationCreateRequest, authenticatedUserId uint) (models.Organization, error)
	Update(organizationId, authenticatedUserId uint, request *requests.OrganizationCreateRequest) (*models.Organization, error)
	Delete(organizationId, authenticatedUserId uint) (bool, error)
}

type OrganizationServiceImpl struct {
	organizationPolicy policies.OrganizationPolicy
	organizationRepo   *repositories.OrganizationRepository
}

func NewOrganizationService(injector *do.Injector) (OrganizationService, error) {
	policy := policies.NewOrganizationPolicy()
	repo := do.MustInvoke[*repositories.OrganizationRepository](injector)

	return &OrganizationServiceImpl{
		organizationPolicy: policy,
		organizationRepo:   repo,
	}, nil
}

func (s *OrganizationServiceImpl) List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authenticatedUserId)
}

func (s *OrganizationServiceImpl) GetByID(id, authenticatedUserId uint) (models.Organization, error) {
	return s.organizationRepo.GetByIDForUser(id, authenticatedUserId)
}

func (s *OrganizationServiceImpl) Create(request *requests.OrganizationCreateRequest, authenticatedUserId uint) (models.Organization, error) {
	if !s.organizationPolicy.CanCreate(authenticatedUserId) {
		return models.Organization{}, errs.NewForbiddenError("organization.error.createForbidden")
	}

	organization := models.Organization{
		Name: request.Name,
	}

	_, err := s.organizationRepo.Create(&organization, authenticatedUserId)
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) Update(organizationId, authenticatedUserId uint, request *requests.OrganizationCreateRequest) (*models.Organization, error) {
	organization, err := s.organizationRepo.GetByIDForUser(organizationId, authenticatedUserId)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organization.ID, authenticatedUserId) {
		return &models.Organization{}, errs.NewForbiddenError("organization.error.updateForbidden")
	}

	err = utils.PopulateModel(&organization, request)
	if err != nil {
		return nil, err
	}

	return s.organizationRepo.Update(organizationId, &organization)
}

func (s *OrganizationServiceImpl) Delete(organizationId, authenticatedUserId uint) (bool, error) {
	organization, err := s.organizationRepo.GetByIDForUser(organizationId, authenticatedUserId)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organization.ID, authenticatedUserId) {
		return false, errs.NewForbiddenError("organization.error.deleteForbidden")
	}

	return s.organizationRepo.Delete(organizationId)
}
