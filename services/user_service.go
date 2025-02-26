package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
	"os"
	"time"
)

type UserService interface {
	Login(request *requests.UserLoginRequest) (models.User, string, error)
	List(paginationParams utils.PaginationParams) ([]models.User, error)
	ExistsByUUID(id uuid.UUID) error
	GetByID(id uuid.UUID) (models.User, error)
	Create(request *requests.UserCreateRequest) (models.User, error)
	Update(userUUID, authUserUUID uuid.UUID, request *requests.UserUpdateRequest) (*models.User, error)
	Delete(userUUID uuid.UUID) (bool, error)
	Logout(userUUID uuid.UUID) error
}

type UserServiceImpl struct {
	userRepo *repositories.UserRepository
}

func NewUserService(injector *do.Injector) (UserService, error) {
	repo := do.MustInvoke[*repositories.UserRepository](injector)

	return &UserServiceImpl{userRepo: repo}, nil
}

func (s *UserServiceImpl) Login(request *requests.UserLoginRequest) (models.User, string, error) {
	user, err := s.userRepo.GetByEmail(request.Email)
	if err != nil {
		return models.User{}, "", err
	}

	if !utils.ComparePassword(user.Password, request.Password) {
		return models.User{}, "", errs.NewUnauthorizedError("user.error.invalidCredentials")
	}

	jwtVersion, err := s.userRepo.CreateJWTVersion(user.Uuid)
	if err != nil {
		return models.User{}, "", err
	}

	token, err := s.generateToken(&user, jwtVersion)
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}

func (s *UserServiceImpl) List(paginationParams utils.PaginationParams) ([]models.User, error) {
	return s.userRepo.List(paginationParams)
}

func (s *UserServiceImpl) GetByID(id uuid.UUID) (models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserServiceImpl) ExistsByUUID(id uuid.UUID) error {
	exists, err := s.userRepo.ExistsByID(id)
	if err != nil {
		return err
	}

	if !exists {
		return errs.NewNotFoundError("user.error.notFound")
	}

	return nil
}

func (s *UserServiceImpl) Create(request *requests.UserCreateRequest) (models.User, error) {
	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Status:   models.UserStatusActive,
		RoleID:   models.UserRoleDeveloper,
	}

	_, err := s.userRepo.Create(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) Update(userUUID, authUserUUID uuid.UUID, request *requests.UserUpdateRequest) (*models.User, error) {
	if !policies.CanUpdateUser(userUUID, authUserUUID) {
		return nil, errs.NewForbiddenError("user.error.updateForbidden")
	}

	user, err := s.userRepo.GetByID(userUUID)
	if err != nil {
		return nil, err
	}

	err = utils.PopulateModel(&user, request)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Update(userUUID, &user)
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

func (s *UserServiceImpl) generateToken(user *models.User, jwtVersion int) (string, error) {
	claims := jwt.MapClaims{
		"version": jwtVersion,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"uuid":    user.Uuid.String(),
		"role_id": user.RoleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
