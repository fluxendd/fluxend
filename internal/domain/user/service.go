package user

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/user"
	"fluxton/pkg/auth"
	"fluxton/pkg/errors"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"os"
	"strings"
	"time"
)

type UserService interface {
	Login(request *user.LoginRequest) (User, string, error)
	List(paginationParams dto.PaginationParams) ([]User, error)
	ExistsByUUID(id uuid.UUID) error
	GetByID(id uuid.UUID) (User, error)
	Create(ctx echo.Context, request *user.CreateRequest) (User, string, error)
	Update(userUUID, authUserUUID uuid.UUID, request *user.UpdateRequest) (*User, error)
	Delete(userUUID uuid.UUID) (bool, error)
	Logout(userUUID uuid.UUID) error
}

type UserServiceImpl struct {
	settingService services.SettingService
	userRepo       *repositories.UserRepository
}

func NewUserService(injector *do.Injector) (UserService, error) {
	settingService := do.MustInvoke[services.SettingService](injector)
	repo := do.MustInvoke[*repositories.UserRepository](injector)

	return &UserServiceImpl{
		settingService: settingService,
		userRepo:       repo,
	}, nil
}

func (s *UserServiceImpl) Login(request *user.LoginRequest) (User, string, error) {
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

func (s *UserServiceImpl) List(paginationParams dto.PaginationParams) ([]User, error) {
	return s.userRepo.List(paginationParams)
}

func (s *UserServiceImpl) GetByID(id uuid.UUID) (User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserServiceImpl) ExistsByUUID(id uuid.UUID) error {
	exists, err := s.userRepo.ExistsByID(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("user.error.notFound")
	}

	return nil
}

func (s *UserServiceImpl) Create(ctx echo.Context, request *user.CreateRequest) (User, string, error) {
	if !s.settingService.GetBool(ctx, "allowRegistrations") {
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
		Status:   UserStatusActive,
		RoleID:   UserRoleOwner,
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

func (s *UserServiceImpl) Update(userUUID, authUserUUID uuid.UUID, request *user.UpdateRequest) (*User, error) {
	if !policies.CanUpdateUser(userUUID, authUserUUID) {
		return nil, errors.NewForbiddenError("user.error.updateForbidden")
	}

	updatedUser, err := s.userRepo.GetByID(userUUID)
	if err != nil {
		return nil, err
	}

	// TODO: COME_BACK_FOR_ME
	/*err = user.PopulateModel(&user, request)
	if err != nil {
		return nil, err
	}*/

	return s.userRepo.Update(userUUID, &updatedUser)
}

func (s *UserServiceImpl) Delete(userUUID uuid.UUID) (bool, error) {
	err := s.ExistsByUUID(userUUID)
	if err != nil {
		return false, err
	}

	return s.userRepo.Delete(userUUID)
}

func (s *UserServiceImpl) Logout(userUUID uuid.UUID) error {
	err := s.ExistsByUUID(userUUID)
	if err != nil {
		return err
	}

	_, err = s.userRepo.CreateJWTVersion(userUUID)

	return err
}

func (s *UserServiceImpl) generateToken(user *User, jwtVersion int) (string, error) {
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
