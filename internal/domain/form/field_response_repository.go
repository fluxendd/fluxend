package form

import (
	"github.com/google/uuid"
)

type FieldResponseRepository interface {
	ListForForm(formUUID uuid.UUID) ([]FormResponse, error)
	GetByUUID(formResponseUUID uuid.UUID) (*FormResponse, error)
	Create(formResponse *FormResponse, formFieldResponse *[]FieldResponse) (*FormResponse, error)
	Delete(formResponseUUID uuid.UUID) error
}
