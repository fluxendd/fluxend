package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	CreatedBy      uuid.UUID `db:"created_by"`
	UpdatedBy      uuid.UUID `db:"updated_by"`
	Name           string    `db:"name"`
	DBName         string    `db:"db_name"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (u Project) GetTableName() string {
	return "projects"
}
