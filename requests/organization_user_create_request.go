package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type OrganizationUserCreateRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r *OrganizationUserCreateRequest) Validate() []string {
	err := validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required.Error("UserID is required")),
	)

	// If no errors, return nil
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
