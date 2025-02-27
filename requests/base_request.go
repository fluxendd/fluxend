package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BaseRequest struct{}

func (r *BaseRequest) ExtractValidationErrors(err error) []string {
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
