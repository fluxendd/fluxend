package user_requests

import (
	"fluxton/internal/api/dto"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateRequest struct {
	dto.BaseRequest
	Bio string `json:"bio"`
}

func (r *UpdateRequest) Validate() []string {
	err := validation.ValidateStruct(r,
		// Bio: optional, length 0-500
		validation.Field(&r.Bio,
			validation.Length(0, 500).Error("Bio must be between 0 and 500 characters"),
		),
	)

	return r.ExtractValidationErrors(err)
}
