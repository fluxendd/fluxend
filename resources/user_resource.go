package resources

import (
	"fluxton/models"
	"github.com/google/uuid"
)

type UserResponse struct {
	Uuid      uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	RoleID    int       `json:"roleId"`
	Bio       string    `json:"bio"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

func UserResource(user *models.User) UserResponse {
	return UserResponse{
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

func UserResourceCollection(users []models.User) []UserResponse {
	resourceUsers := make([]UserResponse, len(users))
	for i, user := range users {
		resourceUsers[i] = UserResource(&user)
	}

	return resourceUsers
}
