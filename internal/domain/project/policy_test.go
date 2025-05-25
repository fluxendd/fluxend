package project

import (
	"errors"
	"fluxend/internal/config/constants"
	"fluxend/internal/domain/auth"
	"fluxend/tests/fixtures/mocks/organization"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolicy_CanCreate_Suite(t *testing.T) {
	t.Run("CanCreate: valid developer user in organization", func(t *testing.T) {
		policy, mockRepo := getTestPolicy(t)

		orgUUID := uuid.New()
		authUser := auth.User{Uuid: uuid.New(), RoleID: constants.UserRoleDeveloper}

		mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(true, nil)

		result := policy.CanCreate(orgUUID, authUser)

		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CanCreate: valid admin user in organization", func(t *testing.T) {
		policy, mockRepo := getTestPolicy(t)

		orgUUID := uuid.New()
		authUser := auth.User{Uuid: uuid.New(), RoleID: constants.UserRoleAdmin}

		mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(true, nil)

		result := policy.CanCreate(orgUUID, authUser)

		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CanCreate: invalid cases", func(t *testing.T) {
		tests := []struct {
			name                 string
			userRole             int
			isOrganizationMember bool
			repositoryError      error
			expectedResult       bool
			expectRepositoryCall bool
		}{
			{
				name:                 "User role below developer",
				userRole:             constants.UserRoleExplorer,
				isOrganizationMember: true,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: false,
			},
			{
				name:                 "Developer user not in organization",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Admin user not in organization",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Repository error for developer",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      errors.New("database error"),
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Repository error for admin",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      errors.New("connection timeout"),
				expectedResult:       false,
				expectRepositoryCall: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				policy, mockRepo := getTestPolicy(t)

				orgUUID := uuid.New()
				authUser := auth.User{
					Uuid:   uuid.New(),
					RoleID: tc.userRole,
				}

				if tc.expectRepositoryCall {
					mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(tc.isOrganizationMember, tc.repositoryError)
				}

				result := policy.CanCreate(orgUUID, authUser)

				assert.Equal(t, tc.expectedResult, result)
				mockRepo.AssertExpectations(t)
			})
		}
	})
}

func TestPolicy_CanAccess_Suite(t *testing.T) {
	t.Run("CanAccess: valid user in organization", func(t *testing.T) {
		roles := []int{constants.UserRoleExplorer, constants.UserRoleDeveloper, constants.UserRoleAdmin}

		for _, role := range roles {
			t.Run("Role: "+string(rune(role)), func(t *testing.T) {
				policy, mockRepo := getTestPolicy(t)

				orgUUID := uuid.New()
				authUser := auth.User{
					Uuid:   uuid.New(),
					RoleID: role,
				}

				mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(true, nil)

				result := policy.CanAccess(orgUUID, authUser)

				assert.True(t, result)
				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("CanAccess: invalid cases", func(t *testing.T) {
		tests := []struct {
			name                 string
			userRole             int
			isOrganizationMember bool
			repositoryError      error
			expectedResult       bool
		}{
			{
				name:                 "User not in organization - Viewer",
				userRole:             constants.UserRoleExplorer,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
			},
			{
				name:                 "User not in organization - Developer",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
			},
			{
				name:                 "User not in organization - Admin",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
			},
			{
				name:                 "Repository error - Viewer",
				userRole:             constants.UserRoleExplorer,
				isOrganizationMember: false,
				repositoryError:      errors.New("database error"),
				expectedResult:       false,
			},
			{
				name:                 "Repository error - Developer",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      errors.New("network timeout"),
				expectedResult:       false,
			},
			{
				name:                 "Repository error - Admin",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      errors.New("connection failed"),
				expectedResult:       false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				policy, mockRepo := getTestPolicy(t)

				orgUUID := uuid.New()
				authUser := auth.User{
					Uuid:   uuid.New(),
					RoleID: tc.userRole,
				}

				mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(tc.isOrganizationMember, tc.repositoryError)

				result := policy.CanAccess(orgUUID, authUser)

				assert.Equal(t, tc.expectedResult, result)
				mockRepo.AssertExpectations(t)
			})
		}
	})
}

func TestPolicy_CanUpdate_Suite(t *testing.T) {
	t.Run("CanUpdate: valid developer user in organization", func(t *testing.T) {
		policy, mockRepo := getTestPolicy(t)

		orgUUID := uuid.New()
		authUser := auth.User{
			Uuid:   uuid.New(),
			RoleID: constants.UserRoleDeveloper,
		}

		mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(true, nil)

		result := policy.CanUpdate(orgUUID, authUser)

		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CanUpdate: valid admin user in organization", func(t *testing.T) {
		policy, mockRepo := getTestPolicy(t)

		orgUUID := uuid.New()
		authUser := auth.User{
			Uuid:   uuid.New(),
			RoleID: constants.UserRoleAdmin,
		}

		mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(true, nil)

		result := policy.CanUpdate(orgUUID, authUser)

		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CanUpdate: invalid cases", func(t *testing.T) {
		tests := []struct {
			name                 string
			userRole             int
			isOrganizationMember bool
			repositoryError      error
			expectedResult       bool
			expectRepositoryCall bool
		}{
			{
				name:                 "User role below developer",
				userRole:             constants.UserRoleExplorer,
				isOrganizationMember: true,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: false,
			},
			{
				name:                 "Developer user not in organization",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Admin user not in organization",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      nil,
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Repository error for developer",
				userRole:             constants.UserRoleDeveloper,
				isOrganizationMember: false,
				repositoryError:      errors.New("database connection lost"),
				expectedResult:       false,
				expectRepositoryCall: true,
			},
			{
				name:                 "Repository error for admin",
				userRole:             constants.UserRoleAdmin,
				isOrganizationMember: false,
				repositoryError:      errors.New("query timeout"),
				expectedResult:       false,
				expectRepositoryCall: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				policy, mockRepo := getTestPolicy(t)

				orgUUID := uuid.New()
				authUser := auth.User{
					Uuid:   uuid.New(),
					RoleID: tc.userRole,
				}

				if tc.expectRepositoryCall {
					mockRepo.On("IsOrganizationMember", orgUUID, authUser.Uuid).Return(tc.isOrganizationMember, tc.repositoryError)
				}

				result := policy.CanUpdate(orgUUID, authUser)

				assert.Equal(t, tc.expectedResult, result)
				mockRepo.AssertExpectations(t)
			})
		}
	})
}

func getTestPolicy(t *testing.T) (*Policy, *organization.MockRepository) {
	mockRepo := organization.NewMockRepository(t)
	policy := &Policy{organizationRepo: mockRepo}

	return policy, mockRepo
}
