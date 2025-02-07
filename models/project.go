package models

import (
	"github.com/google/uuid"
	"strings"
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

func (u Project) GetColumns() string {
	return "id, organization_id, created_by, updated_by, name, db_name, created_at, updated_at"
}

func (u Project) GetColumnsWithAlias(alias string) string {
	columns := strings.Split(u.GetColumns(), ", ")
	for i, field := range columns {
		columns[i] = alias + "." + field
	}

	return strings.Join(columns, ", ")
}
