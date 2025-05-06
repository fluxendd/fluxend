package models

import (
	"github.com/google/uuid"
	"time"
)

type Form struct {
	BaseModel
	Uuid        uuid.UUID `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID `db:"project_uuid" json:"projectUuid"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedBy   uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy   uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

func (u Form) GetTableName() string {
	return "fluxton.forms"
}
