package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type NoteCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (r *NoteCreateRequest) Validate() []string {
	err := validation.ValidateStruct(r,
		validation.Field(&r.Title, validation.Required.Error("Title is required"), validation.Length(3, 100).Error("Title must be between 3 and 100 characters")),
		validation.Field(&r.Content, validation.Required.Error("Content is required"), validation.Length(5, 0).Error("Content must be at least 5 characters")),
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
