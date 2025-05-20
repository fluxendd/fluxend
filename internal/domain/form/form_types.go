package form

import (
	"github.com/google/uuid"
)

type CreateFormInput struct {
	ProjectUUID uuid.UUID `json:"projectUuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (u Form) GetTableName() string {
	return "fluxton.forms"
}
