package organization

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/organization"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/user"
	"fluxton/pkg/errors"
	"github.com/google/uuid"
	"github.com/samber/do"
	"time"
)

type Service interface {
	List(paginationParams dto.PaginationParams, authUserId uuid.UUID) ([]Organization, error)
	GetByID(organizationUUID uuid.UUID, authUser auth.AuthUser) (Organization, error)
	Create(request *organization.CreateRequest, authUser auth.AuthUser) (Organization, error)
	Update(organizationUUID uuid.UUID, authUser auth.AuthUser, request *organization.CreateRequest) (*Organization, error)
	Delete(organizationUUID uuid.UUID, authUser auth.AuthUser) (bool, error)
	ListUsers(organizationUUID uuid.UUID, authUser auth.AuthUser) ([]user.User, error)
	CreateUser(request *organization.MemberCreateRequest, organizationUUID uuid.UUID, authUser auth.AuthUser) (user.User, error)
	DeleteUser(organizationUUID, userID uuid.UUID, authUser auth.AuthUser) error
}

type ServiceImpl struct {
	organizationPolicy *Policy
	organizationRepo   *Repository
	userRepo           *user.Repository
}

func NewOrganizationService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*Policy](injector)
	organizationRepo := do.MustInvoke[*Repository](injector)
	userRepo := do.MustInvoke[*user.Repository](injector)

	return &ServiceImpl{
		organizationPolicy: policy,
		organizationRepo:   organizationRepo,
		userRepo:           userRepo,
	}, nil
}

func (s *ServiceImpl) List(paginationParams dto.PaginationParams, authUserId uuid.UUID) ([]Organization, error) {
	return s.organizationRepo.ListForUser(paginationParams, authUserId)
}

func (s *ServiceImpl) GetByID(organizationUUID uuid.UUID, authUser auth.AuthUser) (Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return Organization{}, err
	}

	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return Organization{}, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return organization, nil
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

func (s *ServiceImpl) Create(request *organization.CreateRequest, authUser auth.AuthUser) (Organization, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return Organization{}, errors.NewForbiddenError("organization.error.createForbidden")
	}

	organization := Organization{
		Name:      request.Name,
		CreatedBy: authUser.Uuid,
		UpdatedBy: authUser.Uuid,
	}

	_, err := s.organizationRepo.Create(&organization, authUser.Uuid)
	if err != nil {
		return Organization{}, err
	}

	return organization, nil
}

func (s *ServiceImpl) Update(organizationUUID uuid.UUID, authUser auth.AuthUser, request *organization.CreateRequest) (*Organization, error) {
	organization, err := s.organizationRepo.GetByUUID(organizationUUID)
	if err != nil {
		return nil, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return &Organization{}, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	// TODO: COME_BACK_FOR_ME
	/*err = organization.PopulateModel(&organization, request)
	if err != nil {
		return nil, err
	}*/

	organization.UpdatedBy = authUser.Uuid
	organization.UpdatedAt = time.Now()

	return s.organizationRepo.Update(&organization)
}

func (s *ServiceImpl) Delete(organizationUUID uuid.UUID, authUser auth.AuthUser) (bool, error) {
	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return false, err
	}

	if !s.organizationPolicy.CanUpdate(organizationUUID, authUser) {
		return false, errors.NewForbiddenError("organization.error.updateForbidden")
	}

	return s.organizationRepo.Delete(organizationUUID)
}

func (s *ServiceImpl) ListUsers(organizationUUID uuid.UUID, authUser auth.AuthUser) ([]user.User, error) {
	if !s.organizationPolicy.CanAccess(organizationUUID, authUser) {
		return nil, errors.NewForbiddenError("organization.error.viewForbidden")
	}

	return s.organizationRepo.ListUsers(organizationUUID)
}

func (s *ServiceImpl) CreateUser(request *organization.MemberCreateRequest, organizationUUID uuid.UUID, authUser auth.AuthUser) (user.User, error) {
	if !s.organizationPolicy.CanCreate(authUser) {
		return user.User{}, errors.NewForbiddenError("organization.error.createUserForbidden")
	}

	err := s.ExistsByUUID(organizationUUID)
	if err != nil {
		return user.User{}, err
	}

	exists, err := s.userRepo.ExistsByID(request.UserID)
	if err != nil {
		return user.User{}, err
	}

	if !exists {
		return user.User{}, errors.NewNotFoundError("user.error.notFound")
	}

	userExists, err := s.organizationRepo.IsOrganizationMember(organizationUUID, request.UserID)
	if err != nil {
		return user.User{}, err
	}

	if userExists {
		return user.User{}, errors.NewUnprocessableError("organization.error.userAlreadyExists")
	}

	err = s.organizationRepo.CreateUser(organizationUUID, request.UserID)
	if err != nil {
		return user.User{}, err
	}

	return s.organizationRepo.GetUser(organizationUUID, request.UserID)
}

func (s *ServiceImpl) DeleteUser(organizationUUID, userUUID uuid.UUID, authUser auth.AuthUser) error {
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
