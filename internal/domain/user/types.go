package user

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	Bio string `json:"bio"`
}
