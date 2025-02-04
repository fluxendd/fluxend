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
	GetByID(id uuid.UUID) (models.User, error)
	Create(request *requests.UserCreateRequest) (models.User, error)
	Update(userId, authenticatedUserId uuid.UUID, request *requests.UserUpdateRequest) (*models.User, error)
	Delete(userId uuid.UUID) (bool, error)
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

	token, err := s.GenerateToken(&user)
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

func (s *UserServiceImpl) Create(request *requests.UserCreateRequest) (models.User, error) {
	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Status:   models.UserStatusActive,
	}

	_, err := s.userRepo.Create(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) Update(userId, authenticatedUserId uuid.UUID, request *requests.UserUpdateRequest) (*models.User, error) {
	if !policies.CanUpdateUser(userId, authenticatedUserId) {
		return nil, errs.NewForbiddenError("user.error.updateForbidden")
	}

	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}

	err = utils.PopulateModel(&user, request)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Update(userId, &user)
}

func (s *UserServiceImpl) Delete(userId uuid.UUID) (bool, error) {
	exists, err := s.userRepo.ExistsByID(userId)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errs.NewNotFoundError("user.error.notFound")
	}

	return s.userRepo.Delete(userId)
}

func (s *UserServiceImpl) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"id":      user.ID,
		"role_id": user.RoleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
