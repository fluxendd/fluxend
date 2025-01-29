package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OrganizationCreateRequest struct {
	Name string `json:"name"`
}

func (r *OrganizationCreateRequest) Validate() []string {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Name is required"), validation.Length(3, 100).Error("Title must be between 3 and 100 characters")),
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
