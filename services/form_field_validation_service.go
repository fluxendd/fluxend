package services

import (
	"encoding/json"
	"fluxton/models"
	"fluxton/pkg"
	"fluxton/pkg/errors"
	"fluxton/requests/form_requests"
	"github.com/samber/do"
)

type FormFieldValidationService interface {
	Validate(value string, formField models.FormField) error
	validateNumber(value string, formField models.FormField) error
	validateString(value string, formField models.FormField) error
	validateEmail(value string, formField models.FormField) error
	validateDate(value string, formField models.FormField) error
	validateCheckbox(value string, formField models.FormField) error
	validateSelect(value string, formField models.FormField) error
}

type FormFieldValidationServiceImpl struct{}

func NewFormFieldValidationService(injector *do.Injector) (FormFieldValidationService, error) {
	return &FormFieldValidationServiceImpl{}, nil
}

func (s *FormFieldValidationServiceImpl) Validate(value string, formField models.FormField) error {
	var validationErr error
	if formField.IsRequired && value == "" {
		return errors.NewUnprocessableError("formResponse.error.fieldRequired")
	}

	switch formField.Type {
	case form_requests.FieldTypeNumber:
		validationErr = s.validateNumber(value, formField)
	case form_requests.FieldTypeText:
	case form_requests.FieldTypeTextarea:
		validationErr = s.validateString(value, formField)
	case form_requests.FieldTypeEmail:
		validationErr = s.validateEmail(value, formField)
	case form_requests.FieldTypeCheckbox:
		validationErr = s.validateCheckbox(value, formField)
	case form_requests.FieldTypeRadio:
	case form_requests.FieldTypeSelect:
		validationErr = s.validateSelect(value, formField)
	}

	return validationErr
}

func (s *FormFieldValidationServiceImpl) validateNumber(value string, formField models.FormField) error {
	numericValue, err := pkg.ConvertStringToInt(value)
	if err != nil {
		return errors.NewUnprocessableError("formResponse.error.invalidNumber")
	}

	if formField.MinValue.Valid && numericValue < int(formField.MinValue.Int64) {
		return errors.NewUnprocessableError("formResponse.error.numberTooLow")
	}

	if formField.MaxValue.Valid && numericValue > int(formField.MaxValue.Int64) {
		return errors.NewUnprocessableError("formResponse.error.numberTooHigh")
	}

	return nil
}

func (s *FormFieldValidationServiceImpl) validateString(value string, formField models.FormField) error {
	if formField.MinLength.Valid && len(value) < int(formField.MinLength.Int64) {
		return errors.NewUnprocessableError("formResponse.error.stringTooShort")
	}

	if formField.MaxLength.Valid && len(value) > int(formField.MaxLength.Int64) {
		return errors.NewUnprocessableError("formResponse.error.stringTooLong")
	}

	if formField.Pattern.Valid {
		matched, err := pkg.MatchRegex(value, formField.Pattern.String)
		if err != nil || !matched {
			return errors.NewUnprocessableError("formResponse.error.invalidPattern")
		}
	}

	return nil
}

func (s *FormFieldValidationServiceImpl) validateEmail(value string, formField models.FormField) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := pkg.MatchRegex(value, emailRegex)
	if err != nil || !matched {
		return errors.NewUnprocessableError("formResponse.error.invalidEmail")
	}

	return nil
}

func (s *FormFieldValidationServiceImpl) validateDate(value string, formField models.FormField) error {
	// TODO: Implement date validation logic

	return nil
}

func (s *FormFieldValidationServiceImpl) validateCheckbox(value string, formField models.FormField) error {
	// TODO: Implement checkbox validation logic

	return nil
}

func (s *FormFieldValidationServiceImpl) validateSelect(value string, formField models.FormField) error {
	var options []string
	err := json.Unmarshal([]byte(formField.Options.String), &options)
	if err != nil {
		return errors.NewUnprocessableError("formResponse.error.invalidSelectOptions")
	}

	if formField.Options.Valid {
		for _, option := range options {
			if option == value {
				return nil
			}
		}

		return errors.NewUnprocessableError("formResponse.error.invalidSelectOption")
	}

	return nil
}
