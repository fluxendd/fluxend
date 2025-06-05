package mapper

import (
	userDto "fluxend/internal/api/dto/user"
	userDomain "fluxend/internal/domain/user"
	"github.com/google/uuid"
)

func ToUserResource(user *userDomain.User) userDto.Response {
	return userDto.Response{
		Uuid:      user.Uuid,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		RoleID:    user.RoleID,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToRegisterUserResource(user *userDomain.User, organizationUUID uuid.UUID) userDto.Response {
	userResponse := ToUserResource(user)
	userResponse.OrganizationUuid = &organizationUUID

	return userResponse
}

func ToUserResourceCollection(users []userDomain.User) []userDto.Response {
	resourceUsers := make([]userDto.Response, len(users))
	for i, currentUser := range users {
		resourceUsers[i] = ToUserResource(&currentUser)
	}

	return resourceUsers
}
