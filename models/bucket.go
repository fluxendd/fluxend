package models

import (
	"github.com/google/uuid"
	"time"
)

type Bucket struct {
	Uuid        uuid.UUID `db:"uuid" json:"uuid"`
	ProjectUuid uuid.UUID `db:"project_uuid" json:"projectUuid"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsPublic    bool      `db:"is_public" json:"isPublic"`
	CreatedBy   uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy   uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

func (u Bucket) GetTableName() string {
	return "storage.buckets"
}
