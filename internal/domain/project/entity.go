package project

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type Project struct {
	shared.BaseEntity
	Uuid             uuid.UUID `db:"uuid"`
	OrganizationUuid uuid.UUID `db:"organization_uuid"`
	CreatedBy        uuid.UUID `db:"created_by"`
	UpdatedBy        uuid.UUID `db:"updated_by"`
	Name             string    `db:"name"`
	Status           string    `db:"status"`
	Description      string    `db:"description"`
	DBName           string    `db:"db_name"`
	DBPort           int       `db:"db_port"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func (u Project) GetTableName() string {
	return "projects"
}
