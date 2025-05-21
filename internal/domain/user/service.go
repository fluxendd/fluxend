package user

import (
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/setting"
	"fluxton/internal/domain/shared"
	"fluxton/pkg/auth"
	"fluxton/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"os"
	"strings"
	"time"
)

type Service interface {
	Login(request *LoginUserInput) (User, string, error)
	List(paginationParams shared.PaginationParams) ([]User, error)
	ExistsByUUID(id uuid.UUID) error
	GetByID(id uuid.UUID) (User, error)
	Create(ctx echo.Context, request *CreateUserInput) (User, string, error)
	Update(userUUID, authUserUUID uuid.UUID, request *UpdateUserInput) (*User, error)
	Delete(userUUID uuid.UUID) (bool, error)
	Logout(userUUID uuid.UUID) error
}

type ServiceImpl struct {
	policy         *Policy
	settingService setting.Service
	userRepo       Repository
}

func NewUserService(injector *do.Injector) (Service, error) {
	policy := do.MustInvoke[*Policy](injector)
	settingService := do.MustInvoke[setting.Service](injector)
	repo := do.MustInvoke[Repository](injector)

	return &ServiceImpl{
		policy:         policy,
		settingService: settingService,
		userRepo:       repo,
	}, nil
}

func (s *ServiceImpl) Login(request *LoginUserInput) (User, string, error) {
	fetchedUser, err := s.userRepo.GetByEmail(request.Email)
	if err != nil {
		return User{}, "", err
	}

	if !auth.ComparePassword(fetchedUser.Password, request.Password) {
		return User{}, "", errors.NewUnauthorizedError("user.error.invalidCredentials")
	}

	jwtVersion, err := s.userRepo.CreateJWTVersion(fetchedUser.Uuid)
	if err != nil {
		return User{}, "", err
	}

	token, err := s.generateToken(&fetchedUser, jwtVersion)
	if err != nil {
		return User{}, "", err
	}

	return fetchedUser, token, nil
}

func (s *ServiceImpl) List(paginationParams shared.PaginationParams) ([]User, error) {
	return s.userRepo.List(paginationParams)
}

func (s *ServiceImpl) GetByID(id uuid.UUID) (User, error) {
	return s.userRepo.GetByID(id)
}

func (s *ServiceImpl) ExistsByUUID(id uuid.UUID) error {
	exists, err := s.userRepo.ExistsByID(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("user.error.notFound")
	}

	return nil
}

func (s *ServiceImpl) Create(ctx echo.Context, request *CreateUserInput) (User, string, error) {
	if !s.settingService.GetBool("allowRegistrations") {
		return User{}, "", errors.NewBadRequestError("user.error.registrationDisabled")
	}

	existsByEmail, err := s.userRepo.ExistsByEmail(request.Email)
	if err != nil {
		return User{}, "", err
	}

	if existsByEmail {
		return User{}, "", errors.NewBadRequestError("user.error.emailAlreadyExists")
	}

	existsByUsername, err := s.userRepo.ExistsByUsername(request.Username)
	if err != nil {
		return User{}, "", err
	}

	if existsByUsername {
		return User{}, "", errors.NewBadRequestError("user.error.usernameAlreadyExists")
	}

	userData := User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Status:   constants.UserStatusActive,
		RoleID:   constants.UserRoleOwner,
	}

	_, err = s.userRepo.Create(&userData)
	if err != nil {
		return User{}, "", err
	}

	jwtVersion, err := s.userRepo.CreateJWTVersion(userData.Uuid)
	if err != nil {
		return User{}, "", err
	}

	token, err := s.generateToken(&userData, jwtVersion)
	if err != nil {
		return User{}, "", err
	}

	return userData, token, nil
}

func (s *ServiceImpl) Update(userUUID, authUserUUID uuid.UUID, request *UpdateUserInput) (*User, error) {
	if !s.policy.CanUpdateUser(userUUID, authUserUUID) {
		return nil, errors.NewForbiddenError("user.error.updateForbidden")
	}

	fetchedUser, err := s.userRepo.GetByID(userUUID)
	if err != nil {
		return nil, err
	}

	err = fetchedUser.PopulateModel(&fetchedUser, request)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Update(userUUID, &fetchedUser)
}

func (s *ServiceImpl) Delete(userUUID uuid.UUID) (bool, error) {
	err := s.ExistsByUUID(userUUID)
	if err != nil {
		return false, err
	}

	return s.userRepo.Delete(userUUID)
}

func (s *ServiceImpl) Logout(userUUID uuid.UUID) error {
	err := s.ExistsByUUID(userUUID)
	if err != nil {
		return err
	}

	_, err = s.userRepo.CreateJWTVersion(userUUID)

	return err
}

func (s *ServiceImpl) generateToken(user *User, jwtVersion int) (string, error) {
	claims := jwt.MapClaims{
		"version": jwtVersion,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"uuid":    user.Uuid.String(),
		"role_id": user.RoleID,                                               // fluxton role
		"role":    "usr_" + strings.ReplaceAll(user.Uuid.String(), "-", "_"), // postgrest role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
