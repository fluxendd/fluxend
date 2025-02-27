package user_requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserLoginRequest) Validate() []string {
	err := validation.ValidateStruct(r,
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
