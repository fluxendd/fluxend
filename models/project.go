package models

import (
	"strings"
	"time"
)

type Project struct {
	ID             uint      `db:"id"`
	OrganizationID uint      `db:"organization_id"`
	Name           string    `db:"name"`
	DBName         string    `db:"db_name"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (u Project) GetTableName() string {
	return "projects"
}

func (u Project) GetFields() string {
	return "id, name, db_name, created_at, updated_at"
}

func (u Project) GetFieldsWithAlias(alias string) string {
	fields := strings.Split(u.GetFields(), ", ")
	for i, field := range fields {
		fields[i] = alias + "." + field
	}

	return strings.Join(fields, ", ")
}
