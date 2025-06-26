package organization

import (
	"fluxend/internal/domain/auth"
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/user"
	"fluxend/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(paginationParams shared.PaginationParams, authUserId uuid.UUID) ([]Organization, error)
	GetByID(organizationUUID uuid.UUID, authUser auth.User) (Organization, error)
	Create(name string, authUser auth.User) (Organization, error)
	Update(name string, organizationUUID uuid.UUID, authUser auth.User) (*Organization, error)
	Delete(organizationUUID uuid.UUID, authUser auth.User) (bool, error)
	ListUsers(organizationUUID uuid.UUID, authUser auth.User) ([]user.User, error)
	CreateUser(userUUID uuid.UUID, organizationUUID uuid.UUID, authUser auth.User) (user.User, error)
	DeleteUser(organizationUUID, userID uuid.UUID, authUser auth.User) error
}

type ServiceImpl struct {
	organizationPolicy *Policy
	organizationRepo   Repository
	userRepo           user.Repository
}

func NewOrganizationService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*Policy](injector)
	organizationRepo := do.MustInvoke[Repository](injector)
	userRepo := do.MustInvoke[user.Repository](injector)

	return &ServiceImpl{
		organizationPolicy: policy,
		organizationRepo:   organizationRepo,
		userRepo:           userRepo,
	}, nil
}

func (s *ServiceImpl) List(paginationParams shared.PaginationParams, authUserId uuid.UUID) ([]Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authUserId)
}

func (s *ServiceImpl) GetByID(organizationUUID uuid.UUID, authUser auth.User) (Organization, error) {
	fetchedOrganization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return Organization{}, err
	}

	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return Organization{}, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return fetchedOrganization, nil
}

func (s *ServiceImpl) ExistsByUUID(organizationUUID uuid.UUID) error {
	exists, err := s.organizationRepo.ExistsByID(organizationUUID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("organization.error.notFound")
	}

	return nil
}

func (s *ServiceImpl) Create(name string, authUser auth.User) (Organization, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return Organization{}, errors.NewForbiddenError("organization.error.createForbidden")
	}

	organizationInput := Organization{
		Name:      name,
		CreatedBy: authUser.Uuid,
		UpdatedBy: authUser.Uuid,
	}

	_, err := s.organizationRepo.Create(&organizationInput, authUser.Uuid)
	if err != nil {
		return Organization{}, err
	}

	return organizationInput, nil
}

func (s *ServiceImpl) Update(name string, organizationUUID uuid.UUID, authUser auth.User) (*Organization, error) {
	fetchedOrganization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return &Organization{}, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	fetchedOrganization.Name = name
	fetchedOrganization.UpdatedBy = authUser.Uuid
	fetchedOrganization.UpdatedAt = time.Now()

	return s.organizationRepo.Update(&fetchedOrganization)
}

func (s *ServiceImpl) Delete(organizationUUID uuid.UUID, authUser auth.User) (bool, error) {
	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	return s.organizationRepo.Delete(organizationUUID)
}

func (s *ServiceImpl) ListUsers(organizationUUID uuid.UUID, authUser auth.User) ([]user.User, error) {
	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return s.organizationRepo.ListUsers(organizationUUID)
}

func (s *ServiceImpl) CreateUser(userUUID uuid.UUID, organizationUUID uuid.UUID, authUser auth.User) (user.User, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return user.User{}, errors.NewForbiddenError("organization.error.createUserForbidden")
	}

	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return user.User{}, err
	}

	exists, err := s.userRepo.ExistsByID(userUUID)
	if err != nil {
		return user.User{}, err
	}

	if !exists {
		return user.User{}, errors.NewNotFoundError("user.error.notFound")
	}

	userExists, err := s.organizationRepo.IsOrganizationMember(organizationUUID, userUUID)
	if err != nil {
		return user.User{}, err
	}

	if userExists {
		return user.User{}, errors.NewUnprocessableError("organization.error.userAlreadyExists")
	}

	if err = s.organizationRepo.CreateUser(organizationUUID, userUUID); err != nil {
		return user.User{}, err
	}

	return s.organizationRepo.GetUser(organizationUUID, userUUID)
}

func (s *ServiceImpl) DeleteUser(organizationUUID, userUUID uuid.UUID, authUser auth.User) error {
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
