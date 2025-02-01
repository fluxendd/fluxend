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

func (u Organization) GetColumns() string {
	return "id, name, created_at, updated_at"
}

func (u Organization) GetColumnsWithAlias(alias string) string {
	columns := strings.Split(u.GetColumns(), ", ")
	for i, column := range columns {
		columns[i] = alias + "." + column
	}

	return strings.Join(columns, ", ")
}
