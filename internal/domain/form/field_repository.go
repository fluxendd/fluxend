package form

import (
	"github.com/google/uuid"
)

type FieldRepository interface {
	ListForForm(formUUID uuid.UUID) ([]Field, error)
	GetByUUID(formUUID uuid.UUID) (Field, error)
	ExistsByUUID(formFieldUUID uuid.UUID) (bool, error)
	ExistsByAnyLabelForForm(labels []string, formUUID uuid.UUID) (bool, error)
	ExistsByLabelForForm(label string, formUUID uuid.UUID) (bool, error)
	Create(formField *Field) (*Field, error)
	CreateMany(formFields []Field, formUUID uuid.UUID) ([]Field, error)
	Update(formField *Field) (*Field, error)
	Delete(formFieldUUID uuid.UUID) (bool, error)
}
