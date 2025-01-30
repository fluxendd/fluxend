package factories

import (
	"github.com/samber/do"
	"myapp/utils"
	"time"

	"myapp/models"
	"myapp/repositories"
)

const defaultPassword = "password"

type UserOption func(*models.User)

type UserFactory struct {
	repo *repositories.UserRepository
}

func NewUserFactory(injector *do.Injector) (*UserFactory, error) {
	repo := do.MustInvoke[*repositories.UserRepository](injector)

	return &UserFactory{repo: repo}, nil
}

// Create a user with options
func (f *UserFactory) Create(opts ...UserOption) (*models.User, error) {
	user := &models.User{
		Username:  utils.Faker.Internet().User(),
		Email:     utils.Faker.Internet().Email(),
		Password:  defaultPassword,
		RoleID:    models.UserRoleBishop,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(user)
	}

	createdUser, err := f.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (f *UserFactory) CreateMany(count int, opts ...UserOption) ([]*models.User, error) {
	var users []*models.User
	for i := 0; i < count; i++ {
		user, err := f.Create(opts...)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (f *UserFactory) WithRole(role int) UserOption {
	return func(user *models.User) {
		user.RoleID = role
	}
}

func (f *UserFactory) WithUsername(username string) UserOption {
	return func(user *models.User) {
		user.Username = username
	}
}

func (f *UserFactory) WithEmail(email string) UserOption {
	return func(user *models.User) {
		user.Email = email
	}
}
