package models

import (
	"github.com/google/uuid"
	"time"
)

type FormField struct {
	Uuid        uuid.UUID `db:"uuid" json:"uuid"`
	FormUuid    uuid.UUID `db:"form_uuid" json:"formUuid"`
	Label       string    `db:"label" json:"label"`
	Description string    `db:"description" json:"description"`
	Type        string    `db:"type" json:"type"`
	IsRequired  bool      `db:"is_required" json:"isRequired"`
	Options     string    `db:"options" json:"options"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

func (ff FormField) GetTableName() string {
	return "fluxton.form_fields"
}
