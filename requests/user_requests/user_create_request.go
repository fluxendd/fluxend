package user_requests

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (r *UserCreateRequest) Validate() []string {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	err := validation.ValidateStruct(r,
		// Username: required, length 3-100, no spaces or special characters
		validation.Field(&r.Username,
			validation.Required.Error("Username is required"),
			validation.Length(3, 100).Error("Username must be between 3 and 100 characters"),
			validation.Match(usernameRegex).Error("Username must not contain spaces or special characters"),
		),
		// Email: required, valid format
		validation.Field(&r.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Email must be a valid email address"),
		),
		// Password: required, at least 5 characters
		validation.Field(&r.Password,
			validation.Required.Error("Password is required"),
			validation.Length(5, 0).Error("Password must be at least 5 characters"),
		),
		// Bio: optional, length 0-500
		validation.Field(&r.Bio,
			validation.Length(0, 500).Error("Bio must be between 0 and 500 characters"),
		),
	)

	if err == nil {
		return nil
	}

	var errors []string
	if ve, ok := err.(validation.Errors); ok {
		for _, validationErr := range ve {
			errors = append(errors, validationErr.Error())
		}
	}

	return errors
}
