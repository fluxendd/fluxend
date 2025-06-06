// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package organization

import (
	"fluxend/internal/domain/organization"
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/user"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

type MockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepository) EXPECT() *MockRepository_Expecter {
	return &MockRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function for the type MockRepository
func (_mock *MockRepository) Create(organization1 *organization.Organization, authUserID uuid.UUID) (*organization.Organization, error) {
	ret := _mock.Called(organization1, authUserID)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *organization.Organization
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(*organization.Organization, uuid.UUID) (*organization.Organization, error)); ok {
		return returnFunc(organization1, authUserID)
	}
	if returnFunc, ok := ret.Get(0).(func(*organization.Organization, uuid.UUID) *organization.Organization); ok {
		r0 = returnFunc(organization1, authUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*organization.Organization)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*organization.Organization, uuid.UUID) error); ok {
		r1 = returnFunc(organization1, authUserID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - organization1
//   - authUserID
func (_e *MockRepository_Expecter) Create(organization1 interface{}, authUserID interface{}) *MockRepository_Create_Call {
	return &MockRepository_Create_Call{Call: _e.mock.On("Create", organization1, authUserID)}
}

func (_c *MockRepository_Create_Call) Run(run func(organization1 *organization.Organization, authUserID uuid.UUID)) *MockRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*organization.Organization), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_Create_Call) Return(organization11 *organization.Organization, err error) *MockRepository_Create_Call {
	_c.Call.Return(organization11, err)
	return _c
}

func (_c *MockRepository_Create_Call) RunAndReturn(run func(organization1 *organization.Organization, authUserID uuid.UUID) (*organization.Organization, error)) *MockRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function for the type MockRepository
func (_mock *MockRepository) CreateUser(organizationUUID uuid.UUID, userUUID uuid.UUID) error {
	ret := _mock.Called(organizationUUID, userUUID)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = returnFunc(organizationUUID, userUUID)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRepository_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockRepository_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - organizationUUID
//   - userUUID
func (_e *MockRepository_Expecter) CreateUser(organizationUUID interface{}, userUUID interface{}) *MockRepository_CreateUser_Call {
	return &MockRepository_CreateUser_Call{Call: _e.mock.On("CreateUser", organizationUUID, userUUID)}
}

func (_c *MockRepository_CreateUser_Call) Run(run func(organizationUUID uuid.UUID, userUUID uuid.UUID)) *MockRepository_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_CreateUser_Call) Return(err error) *MockRepository_CreateUser_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRepository_CreateUser_Call) RunAndReturn(run func(organizationUUID uuid.UUID, userUUID uuid.UUID) error) *MockRepository_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function for the type MockRepository
func (_mock *MockRepository) Delete(organizationUUID uuid.UUID) (bool, error) {
	ret := _mock.Called(organizationUUID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) (bool, error)); ok {
		return returnFunc(organizationUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = returnFunc(organizationUUID)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - organizationUUID
func (_e *MockRepository_Expecter) Delete(organizationUUID interface{}) *MockRepository_Delete_Call {
	return &MockRepository_Delete_Call{Call: _e.mock.On("Delete", organizationUUID)}
}

func (_c *MockRepository_Delete_Call) Run(run func(organizationUUID uuid.UUID)) *MockRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_Delete_Call) Return(b bool, err error) *MockRepository_Delete_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *MockRepository_Delete_Call) RunAndReturn(run func(organizationUUID uuid.UUID) (bool, error)) *MockRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUser provides a mock function for the type MockRepository
func (_mock *MockRepository) DeleteUser(organizationUUID uuid.UUID, userUUID uuid.UUID) error {
	ret := _mock.Called(organizationUUID, userUUID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = returnFunc(organizationUUID, userUUID)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRepository_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type MockRepository_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - organizationUUID
//   - userUUID
func (_e *MockRepository_Expecter) DeleteUser(organizationUUID interface{}, userUUID interface{}) *MockRepository_DeleteUser_Call {
	return &MockRepository_DeleteUser_Call{Call: _e.mock.On("DeleteUser", organizationUUID, userUUID)}
}

func (_c *MockRepository_DeleteUser_Call) Run(run func(organizationUUID uuid.UUID, userUUID uuid.UUID)) *MockRepository_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_DeleteUser_Call) Return(err error) *MockRepository_DeleteUser_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRepository_DeleteUser_Call) RunAndReturn(run func(organizationUUID uuid.UUID, userUUID uuid.UUID) error) *MockRepository_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsByID provides a mock function for the type MockRepository
func (_mock *MockRepository) ExistsByID(organizationUUID uuid.UUID) (bool, error) {
	ret := _mock.Called(organizationUUID)

	if len(ret) == 0 {
		panic("no return value specified for ExistsByID")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) (bool, error)); ok {
		return returnFunc(organizationUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = returnFunc(organizationUUID)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_ExistsByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsByID'
type MockRepository_ExistsByID_Call struct {
	*mock.Call
}

// ExistsByID is a helper method to define mock.On call
//   - organizationUUID
func (_e *MockRepository_Expecter) ExistsByID(organizationUUID interface{}) *MockRepository_ExistsByID_Call {
	return &MockRepository_ExistsByID_Call{Call: _e.mock.On("ExistsByID", organizationUUID)}
}

func (_c *MockRepository_ExistsByID_Call) Run(run func(organizationUUID uuid.UUID)) *MockRepository_ExistsByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_ExistsByID_Call) Return(b bool, err error) *MockRepository_ExistsByID_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *MockRepository_ExistsByID_Call) RunAndReturn(run func(organizationUUID uuid.UUID) (bool, error)) *MockRepository_ExistsByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUUID provides a mock function for the type MockRepository
func (_mock *MockRepository) GetByUUID(organizationUUID uuid.UUID) (organization.Organization, error) {
	ret := _mock.Called(organizationUUID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUUID")
	}

	var r0 organization.Organization
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) (organization.Organization, error)); ok {
		return returnFunc(organizationUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) organization.Organization); ok {
		r0 = returnFunc(organizationUUID)
	} else {
		r0 = ret.Get(0).(organization.Organization)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_GetByUUID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUUID'
type MockRepository_GetByUUID_Call struct {
	*mock.Call
}

// GetByUUID is a helper method to define mock.On call
//   - organizationUUID
func (_e *MockRepository_Expecter) GetByUUID(organizationUUID interface{}) *MockRepository_GetByUUID_Call {
	return &MockRepository_GetByUUID_Call{Call: _e.mock.On("GetByUUID", organizationUUID)}
}

func (_c *MockRepository_GetByUUID_Call) Run(run func(organizationUUID uuid.UUID)) *MockRepository_GetByUUID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_GetByUUID_Call) Return(organization1 organization.Organization, err error) *MockRepository_GetByUUID_Call {
	_c.Call.Return(organization1, err)
	return _c
}

func (_c *MockRepository_GetByUUID_Call) RunAndReturn(run func(organizationUUID uuid.UUID) (organization.Organization, error)) *MockRepository_GetByUUID_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function for the type MockRepository
func (_mock *MockRepository) GetUser(organizationUUID uuid.UUID, userUUID uuid.UUID) (user.User, error) {
	ret := _mock.Called(organizationUUID, userUUID)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 user.User
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) (user.User, error)); ok {
		return returnFunc(organizationUUID, userUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) user.User); ok {
		r0 = returnFunc(organizationUUID, userUUID)
	} else {
		r0 = ret.Get(0).(user.User)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID, userUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type MockRepository_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - organizationUUID
//   - userUUID
func (_e *MockRepository_Expecter) GetUser(organizationUUID interface{}, userUUID interface{}) *MockRepository_GetUser_Call {
	return &MockRepository_GetUser_Call{Call: _e.mock.On("GetUser", organizationUUID, userUUID)}
}

func (_c *MockRepository_GetUser_Call) Run(run func(organizationUUID uuid.UUID, userUUID uuid.UUID)) *MockRepository_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_GetUser_Call) Return(user1 user.User, err error) *MockRepository_GetUser_Call {
	_c.Call.Return(user1, err)
	return _c
}

func (_c *MockRepository_GetUser_Call) RunAndReturn(run func(organizationUUID uuid.UUID, userUUID uuid.UUID) (user.User, error)) *MockRepository_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// IsOrganizationMember provides a mock function for the type MockRepository
func (_mock *MockRepository) IsOrganizationMember(organizationUUID uuid.UUID, authUserID uuid.UUID) (bool, error) {
	ret := _mock.Called(organizationUUID, authUserID)

	if len(ret) == 0 {
		panic("no return value specified for IsOrganizationMember")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) (bool, error)); ok {
		return returnFunc(organizationUUID, authUserID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) bool); ok {
		r0 = returnFunc(organizationUUID, authUserID)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID, authUserID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_IsOrganizationMember_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsOrganizationMember'
type MockRepository_IsOrganizationMember_Call struct {
	*mock.Call
}

// IsOrganizationMember is a helper method to define mock.On call
//   - organizationUUID
//   - authUserID
func (_e *MockRepository_Expecter) IsOrganizationMember(organizationUUID interface{}, authUserID interface{}) *MockRepository_IsOrganizationMember_Call {
	return &MockRepository_IsOrganizationMember_Call{Call: _e.mock.On("IsOrganizationMember", organizationUUID, authUserID)}
}

func (_c *MockRepository_IsOrganizationMember_Call) Run(run func(organizationUUID uuid.UUID, authUserID uuid.UUID)) *MockRepository_IsOrganizationMember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_IsOrganizationMember_Call) Return(b bool, err error) *MockRepository_IsOrganizationMember_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *MockRepository_IsOrganizationMember_Call) RunAndReturn(run func(organizationUUID uuid.UUID, authUserID uuid.UUID) (bool, error)) *MockRepository_IsOrganizationMember_Call {
	_c.Call.Return(run)
	return _c
}

// ListForUser provides a mock function for the type MockRepository
func (_mock *MockRepository) ListForUser(paginationParams shared.PaginationParams, authUserID uuid.UUID) ([]organization.Organization, error) {
	ret := _mock.Called(paginationParams, authUserID)

	if len(ret) == 0 {
		panic("no return value specified for ListForUser")
	}

	var r0 []organization.Organization
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(shared.PaginationParams, uuid.UUID) ([]organization.Organization, error)); ok {
		return returnFunc(paginationParams, authUserID)
	}
	if returnFunc, ok := ret.Get(0).(func(shared.PaginationParams, uuid.UUID) []organization.Organization); ok {
		r0 = returnFunc(paginationParams, authUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]organization.Organization)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(shared.PaginationParams, uuid.UUID) error); ok {
		r1 = returnFunc(paginationParams, authUserID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_ListForUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListForUser'
type MockRepository_ListForUser_Call struct {
	*mock.Call
}

// ListForUser is a helper method to define mock.On call
//   - paginationParams
//   - authUserID
func (_e *MockRepository_Expecter) ListForUser(paginationParams interface{}, authUserID interface{}) *MockRepository_ListForUser_Call {
	return &MockRepository_ListForUser_Call{Call: _e.mock.On("ListForUser", paginationParams, authUserID)}
}

func (_c *MockRepository_ListForUser_Call) Run(run func(paginationParams shared.PaginationParams, authUserID uuid.UUID)) *MockRepository_ListForUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(shared.PaginationParams), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_ListForUser_Call) Return(organizations []organization.Organization, err error) *MockRepository_ListForUser_Call {
	_c.Call.Return(organizations, err)
	return _c
}

func (_c *MockRepository_ListForUser_Call) RunAndReturn(run func(paginationParams shared.PaginationParams, authUserID uuid.UUID) ([]organization.Organization, error)) *MockRepository_ListForUser_Call {
	_c.Call.Return(run)
	return _c
}

// ListUsers provides a mock function for the type MockRepository
func (_mock *MockRepository) ListUsers(organizationUUID uuid.UUID) ([]user.User, error) {
	ret := _mock.Called(organizationUUID)

	if len(ret) == 0 {
		panic("no return value specified for ListUsers")
	}

	var r0 []user.User
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) ([]user.User, error)); ok {
		return returnFunc(organizationUUID)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) []user.User); ok {
		r0 = returnFunc(organizationUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(organizationUUID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_ListUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsers'
type MockRepository_ListUsers_Call struct {
	*mock.Call
}

// ListUsers is a helper method to define mock.On call
//   - organizationUUID
func (_e *MockRepository_Expecter) ListUsers(organizationUUID interface{}) *MockRepository_ListUsers_Call {
	return &MockRepository_ListUsers_Call{Call: _e.mock.On("ListUsers", organizationUUID)}
}

func (_c *MockRepository_ListUsers_Call) Run(run func(organizationUUID uuid.UUID)) *MockRepository_ListUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockRepository_ListUsers_Call) Return(users []user.User, err error) *MockRepository_ListUsers_Call {
	_c.Call.Return(users, err)
	return _c
}

func (_c *MockRepository_ListUsers_Call) RunAndReturn(run func(organizationUUID uuid.UUID) ([]user.User, error)) *MockRepository_ListUsers_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function for the type MockRepository
func (_mock *MockRepository) Update(organization1 *organization.Organization) (*organization.Organization, error) {
	ret := _mock.Called(organization1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *organization.Organization
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(*organization.Organization) (*organization.Organization, error)); ok {
		return returnFunc(organization1)
	}
	if returnFunc, ok := ret.Get(0).(func(*organization.Organization) *organization.Organization); ok {
		r0 = returnFunc(organization1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*organization.Organization)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*organization.Organization) error); ok {
		r1 = returnFunc(organization1)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - organization1
func (_e *MockRepository_Expecter) Update(organization1 interface{}) *MockRepository_Update_Call {
	return &MockRepository_Update_Call{Call: _e.mock.On("Update", organization1)}
}

func (_c *MockRepository_Update_Call) Run(run func(organization1 *organization.Organization)) *MockRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*organization.Organization))
	})
	return _c
}

func (_c *MockRepository_Update_Call) Return(organization11 *organization.Organization, err error) *MockRepository_Update_Call {
	_c.Call.Return(organization11, err)
	return _c
}

func (_c *MockRepository_Update_Call) RunAndReturn(run func(organization1 *organization.Organization) (*organization.Organization, error)) *MockRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}
