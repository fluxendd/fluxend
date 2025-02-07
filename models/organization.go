package models

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Organization struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedBy uuid.UUID `db:"updated_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u Organization) GetTableName() string {
	return "fluxton.organizations"
}

func (u Organization) GetColumns() string {
	return "id, name, created_by, updated_by, created_at, updated_at"
}

func (u Organization) GetColumnsWithAlias(alias string) string {
	columns := strings.Split(u.GetColumns(), ", ")
	for i, column := range columns {
		columns[i] = alias + "." + column
	}

	return strings.Join(columns, ", ")
}
