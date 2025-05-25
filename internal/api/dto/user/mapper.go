package user

import (
	"fluxend/internal/domain/user"
)

func ToCreateUserInput(request *CreateRequest) *user.CreateUserInput {
	return &user.CreateUserInput{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Bio:      request.Bio,
	}
}

func ToLoginUserInput(request *LoginRequest) *user.LoginUserInput {
	return &user.LoginUserInput{
		Email:    request.Email,
		Password: request.Password,
	}
}

func ToUpdateUserInput(request *UpdateRequest) *user.UpdateUserInput {
	return &user.UpdateUserInput{
		Bio: request.Bio,
	}
}
