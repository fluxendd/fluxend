package user

import (
	dtoUser "fluxton/internal/api/dto/user"
	domainUser "fluxton/internal/domain/user"
)

func ToResponse(user *domainUser.User) dtoUser.Response {
	return dtoUser.Response{
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

func ToResponseCollection(users []domainUser.User) []dtoUser.Response {
	resourceUsers := make([]dtoUser.Response, len(users))
	for i, currentUser := range users {
		resourceUsers[i] = ToResponse(&currentUser)
	}

	return resourceUsers
}
