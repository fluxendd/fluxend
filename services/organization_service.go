package services

import (
	"fluxton/internal/api/dto"
	"fluxton/models"
	"fluxton/pkg/errors"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests/organization_requests"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type OrganizationService interface {
	List(paginationParams dto.PaginationParams, authUserId uuid.UUID) ([]models.Organization, error)
	GetByID(organizationUUID uuid.UUID, authUser models.AuthUser) (models.Organization, error)
	Create(request *organization_requests.CreateRequest, authUser models.AuthUser) (models.Organization, error)
	Update(organizationUUID uuid.UUID, authUser models.AuthUser, request *organization_requests.CreateRequest) (*models.Organization, error)
	Delete(organizationUUID uuid.UUID, authUser models.AuthUser) (bool, error)
	ListUsers(organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.User, error)
	CreateUser(request *organization_requests.MemberCreateRequest, organizationUUID uuid.UUID, authUser models.AuthUser) (models.User, error)
	DeleteUser(organizationUUID, userID uuid.UUID, authUser models.AuthUser) error
}

type OrganizationServiceImpl struct {
	organizationPolicy *policies.OrganizationPolicy
	organizationRepo   *repositories.OrganizationRepository
	userRepo           *repositories.UserRepository
}

func NewOrganizationService(injector *do.Injector) (OrganizationService, error) {
	policy := do.MustInvoke[*policies.OrganizationPolicy](injector)
	organizationRepo := do.MustInvoke[*repositories.OrganizationRepository](injector)
	userRepo := do.MustInvoke[*repositories.UserRepository](injector)

	return &OrganizationServiceImpl{
		organizationPolicy: policy,
		organizationRepo:   organizationRepo,
		userRepo:           userRepo,
	}, nil
}

func (s *OrganizationServiceImpl) List(paginationParams dto.PaginationParams, authUserId uuid.UUID) ([]models.Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authUserId)
}

func (s *OrganizationServiceImpl) GetByID(organizationUUID uuid.UUID, authUser models.AuthUser) (models.Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return models.Organization{}, err
	}

	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return models.Organization{}, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return organization, nil
}

func (s *OrganizationServiceImpl) ExistsByUUID(organizationUUID uuid.UUID) error {
	exists, err := s.organizationRepo.ExistsByID(organizationUUID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("organization.error.notFound")
	}

	return nil
}

func (s *OrganizationServiceImpl) Create(request *organization_requests.CreateRequest, authUser models.AuthUser) (models.Organization, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return models.Organization{}, errors.NewForbiddenError("organization.error.createForbidden")
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

func (s *OrganizationServiceImpl) Update(organizationUUID uuid.UUID, authUser models.AuthUser, request *organization_requests.CreateRequest) (*models.Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return &models.Organization{}, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	err = organization.PopulateModel(&organization, request)
	if err != nil {
		return nil, err
	}

	organization.UpdatedBy = authUser.Uuid
	organization.UpdatedAt = time.Now()

	return s.organizationRepo.Update(&organization)
}

func (s *OrganizationServiceImpl) Delete(organizationUUID uuid.UUID, authUser models.AuthUser) (bool, error) {
	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	return s.organizationRepo.Delete(organizationUUID)
}

func (s *OrganizationServiceImpl) ListUsers(organizationUUID uuid.UUID, authUser models.AuthUser) ([]models.User, error) {
	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return s.organizationRepo.ListUsers(organizationUUID)
}

func (s *OrganizationServiceImpl) CreateUser(request *organization_requests.MemberCreateRequest, organizationUUID uuid.UUID, authUser models.AuthUser) (models.User, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return models.User{}, errors.NewForbiddenError("organization.error.createUserForbidden")
	}

	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return models.User{}, err
	}

	exists, err := s.userRepo.ExistsByID(request.UserID)
	if err != nil {
		return models.User{}, err
	}

	if !exists {
		return models.User{}, errors.NewNotFoundError("user.error.notFound")
	}

	userExists, err := s.organizationRepo.IsOrganizationMember(organizationUUID, request.UserID)
	if err != nil {
		return models.User{}, err
	}

	if userExists {
		return models.User{}, errors.NewUnprocessableError("organization.error.userAlreadyExists")
	}

	err = s.organizationRepo.CreateUser(organizationUUID, request.UserID)
	if err != nil {
		return models.User{}, err
	}

	return s.organizationRepo.GetUser(organizationUUID, request.UserID)
}

func (s *OrganizationServiceImpl) DeleteUser(organizationUUID, userUUID uuid.UUID, authUser models.AuthUser) error {
	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return errors.NewForbiddenError("organization.error.deleteUserForbidden")
	}

	userExists, err := s.organizationRepo.IsOrganizationMember(organizationUUID, userUUID)
	if err != nil {
		return err
	}

	if !userExists {
		return errors.NewNotFoundError("organization.error.userNotFound")
	}

	return s.organizationRepo.DeleteUser(organizationUUID, userUUID)
}
