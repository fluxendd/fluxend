package models

import (
	"strings"
	"time"
)

type Organization struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u Organization) GetTableName() string {
	return "organizations"
}

func (u Organization) GetFields() string {
	return "id, name, created_at, updated_at"
}

func (u Organization) GetFieldsWithAlias(alias string) string {
	fields := strings.Split(u.GetFields(), ", ")
	for i, field := range fields {
		fields[i] = alias + "." + field
	}

	return strings.Join(fields, ", ")
}
