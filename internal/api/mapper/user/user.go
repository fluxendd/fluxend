package user

import (
	userDto "fluxton/internal/api/dto/user"
	userDomain "fluxton/internal/domain/user"
)

func ToResponse(user *userDomain.User) userDto.Response {
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

func ToResponseCollection(users []userDomain.User) []userDto.Response {
	resourceUsers := make([]userDto.Response, len(users))
	for i, currentUser := range users {
		resourceUsers[i] = ToResponse(&currentUser)
	}

	return resourceUsers
}
