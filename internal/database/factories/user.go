package factories

import (
	"fluxend/internal/config/constants"
	"fluxend/internal/domain/organization"
	"fluxend/internal/domain/user"
	"fluxend/pkg"
	"github.com/samber/do"
	"time"
)

const defaultPassword = "password"

type UserOption func(user *user.User)

type UserFactory struct {
	userRepo         user.Repository
	organizationRepo organization.Repository
}

func NewUserFactory(injector *do.Injector) (*UserFactory, error) {
	userRepo := do.MustInvoke[user.Repository](injector)
	organizationRepo := do.MustInvoke[organization.Repository](injector)

	return &UserFactory{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}, nil
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

	createdUser, err := f.userRepo.Create(inputUser)
	if err != nil {
		return nil, err
	}

	inputOrganization := &organization.Organization{
		Name:      pkg.Faker.Company().Name(),
		CreatedBy: createdUser.Uuid,
		UpdatedBy: createdUser.Uuid,
	}

	_, err = f.organizationRepo.Create(inputOrganization, createdUser.Uuid)
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

func (f *UserFactory) WithPassword(password string) UserOption {
	return func(user *user.User) {
		user.Password = password
	}
}
