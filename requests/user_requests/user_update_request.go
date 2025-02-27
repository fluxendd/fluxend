package user_requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UserUpdateRequest struct {
	Bio string `json:"bio"`
}

func (r *UserUpdateRequest) Validate() []string {
	err := validation.ValidateStruct(r,
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
