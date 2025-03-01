package models

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Uuid        uuid.UUID `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID `db:"project_uuid" json:"projectUuid"`
	Name        string    `db:"name" json:"name"`
	Size        int       `db:"size" json:"size"` // in KB
	MimeType    string    `db:"mime_type" json:"mimeType"`
	Path        string    `db:"path" json:"path"`
	CreatedBy   uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy   uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

func (u File) GetTableName() string {
	return "storage.files"
}
