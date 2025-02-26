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
	"time"
)

type OrganizationService interface {
	List(paginationParams utils.PaginationParams, authUserId uuid.UUID) ([]models.Organization, error)
	GetByID(organizationUUID uuid.UUID, authUser models.AuthUser) (models.Organization, error)
	Create(request *requests.OrganizationCreateRequest, authUser models.AuthUser) (models.Organization, error)
	Update(organizationUUID uuid.UUID, authUser models.AuthUser, request *requests.OrganizationCreateRequest) (*models.Organization, error)
	Delete(organizationUUID uuid.UUID, authUser models.AuthUser) (bool, error)
	ListUsers(organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.User, error)
	CreateUser(request *requests.OrganizationUserCreateRequest, organizationUUID uuid.UUID, authUser models.AuthUser) (models.User, error)
	DeleteUser(organizationUUID, userID uuid.UUID, authUser models.AuthUser) error
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

func (s *OrganizationServiceImpl) List(paginationParams utils.PaginationParams, authUserId uuid.UUID) ([]models.Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authUserId)
}

func (s *OrganizationServiceImpl) GetByID(organizationUUID uuid.UUID, authUser models.AuthUser) (models.Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return models.Organization{}, err
	}

	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return models.Organization{}, errs.NewForbiddenError("organization.error.viewForbidden")
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) Create(request *requests.OrganizationCreateRequest, authUser models.AuthUser) (models.Organization, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return models.Organization{}, errs.NewForbiddenError("organization.error.createForbidden")
	}

	organization := models.Organization{
		Name:      request.Name,
		CreatedBy: authUser.Uuid,
		UpdatedBy: authUser.Uuid,
	}

	_, err := s.organizationRepo.Create(&organization, authUser.Uuid)
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) Update(organizationUUID uuid.UUID, authUser models.AuthUser, request *requests.OrganizationCreateRequest) (*models.Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.Organization{}, errs.NewForbiddenError("organization.error.updateForbidden")
	}

	err = utils.PopulateModel(&organization, request)
	if err != nil {
		return nil, err
	}

	organization.UpdatedBy = authUser.Uuid
	organization.UpdatedAt = time.Now()

	return s.organizationRepo.Update(&organization)
}

func (s *OrganizationServiceImpl) Delete(organizationUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	_, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errs.NewForbiddenError("organization.error.updateForbidden")
	}

	return s.organizationRepo.Delete(organizationUUID)
}

func (s *OrganizationServiceImpl) ListUsers(organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.User, error) {
	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errs.NewForbiddenError("organization.error.viewForbidden")
	}

	return s.organizationRepo.ListUsers(organizationUUID)
}

func (s *OrganizationServiceImpl) CreateUser(request *requests.OrganizationUserCreateRequest, organizationUUID uuid.UUID, authUser models.AuthUser) (models.User, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return models.User{}, errs.NewForbiddenError("organization.error.createUserForbidden")
	}

	userExists, err := s.organizationRepo.IsOrganizationMember(organizationUUID, request.UserID)
	if err != nil {
		return models.User{}, err
	}

	if userExists {
		return models.User{}, errs.NewUnprocessableError("organization.error.userAlreadyExists")
	}

	err = s.organizationRepo.CreateUser(organizationUUID, request.UserID)
	if err != nil {
		return models.User{}, err
	}

	return s.organizationRepo.GetUser(organizationUUID, request.UserID)
}

func (s *OrganizationServiceImpl) DeleteUser(organizationUUID, userID uuid.UUID, authUser models.AuthUser) error {
	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return errs.NewForbiddenError("organization.error.deleteUserForbidden")
	}

	return s.organizationRepo.DeleteUser(organizationUUID, userID)
}
