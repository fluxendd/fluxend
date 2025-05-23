package factories

import (
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/user"
	"fluxton/pkg"
	"github.com/samber/do"
	"time"
)

const defaultPassword = "password"

type UserOption func(user *user.User)

type UserFactory struct {
	repo user.Repository
}

func NewUserFactory(injector *do.Injector) (*UserFactory, error) {
	repo := do.MustInvoke[user.Repository](injector)

	return &UserFactory{repo: repo}, nil
}

// Create a user with options
func (f *UserFactory) Create(opts ...UserOption) (*user.User, error) {
	inputUser := &user.User{
		Username:  pkg.Faker.Internet().User(),
		Email:     pkg.Faker.Internet().Email(),
		Password:  defaultPassword,
		RoleID:    constants.UserRoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(inputUser)
	}

	createdUser, err := f.repo.Create(inputUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (f *UserFactory) CreateMany(count int, opts ...UserOption) ([]*user.User, error) {
	var users []*user.User
	for i := 0; i < count; i++ {
		currentUser, err := f.Create(opts...)
		if err != nil {
			return nil, err
		}
		users = append(users, currentUser)
	}
	return users, nil
}

func (f *UserFactory) WithRole(role int) UserOption {
	return func(user *user.User) {
		user.RoleID = role
	}
}

func (f *UserFactory) WithUsername(username string) UserOption {
	return func(user *user.User) {
		user.Username = username
	}
}

func (f *UserFactory) WithEmail(email string) UserOption {
	return func(user *user.User) {
		user.Email = email
	}
}
