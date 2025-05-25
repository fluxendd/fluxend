package form

import (
	"encoding/json"
	"fluxend/internal/config/constants"
	"fluxend/pkg"
	"fluxend/pkg/errors"
	"github.com/samber/do"
)

type FieldValidationService interface {
	Validate(value string, formField Field) error
	validateNumber(value string, formField Field) error
	validateString(value string, formField Field) error
	validateEmail(value string, formField Field) error
	validateDate(value string, formField Field) error
	validateCheckbox(value string, formField Field) error
	validateSelect(value string, formField Field) error
}

type FieldValidationServiceImpl struct{}

func NewFormFieldValidationService(injector *do.Injector) (FieldValidationService, error) {
	return &FieldValidationServiceImpl{}, nil
}

func (s *FieldValidationServiceImpl) Validate(value string, formField Field) error {
	var validationErr error
	if formField.IsRequired && value == "" {
		return errors.NewUnprocessableError("formResponse.error.fieldRequired")
	}

	switch formField.Type {
	case constants.FieldTypeNumber:
		validationErr = s.validateNumber(value, formField)
	case constants.FieldTypeText:
	case constants.FieldTypeTextarea:
		validationErr = s.validateString(value, formField)
	case constants.FieldTypeEmail:
		validationErr = s.validateEmail(value, formField)
	case constants.FieldTypeCheckbox:
		validationErr = s.validateCheckbox(value, formField)
	case constants.FieldTypeRadio:
	case constants.FieldTypeSelect:
		validationErr = s.validateSelect(value, formField)
	}

	return validationErr
}

func (s *FieldValidationServiceImpl) validateNumber(value string, formField Field) error {
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

func (s *FieldValidationServiceImpl) validateString(value string, formField Field) error {
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

func (s *FieldValidationServiceImpl) validateEmail(value string, formField Field) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := pkg.MatchRegex(value, emailRegex)
	if err != nil || !matched {
		return errors.NewUnprocessableError("formResponse.error.invalidEmail")
	}

	return nil
}

func (s *FieldValidationServiceImpl) validateDate(value string, formField Field) error {
	// TODO: Implement date validation logic

	return nil
}

func (s *FieldValidationServiceImpl) validateCheckbox(value string, formField Field) error {
	// TODO: Implement checkbox validation logic

	return nil
}

func (s *FieldValidationServiceImpl) validateSelect(value string, formField Field) error {
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
