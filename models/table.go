package models

import (
	"encoding/json"
	"fluxton/types"
	"strings"
	"time"
)

type Table struct {
	ID        uint               `db:"id"`
	ProjectID uint               `db:"project_id"`
	Name      string             `db:"name"`
	Fields    []types.TableField `db:"fields"`
	CreatedAt time.Time          `db:"created_at"`
	UpdatedAt time.Time          `db:"updated_at"`
}

func (u Table) GetTableName() string {
	return "tables"
}

func (u Table) GetFields() string {
	return "id, project_id, name, fields, created_at, updated_at"
}

func (u Table) GetFieldsWithAlias(alias string) string {
	fields := strings.Split(u.GetFields(), ", ")
	for i, field := range fields {
		fields[i] = alias + "." + field
	}

	return strings.Join(fields, ", ")
}

func (t Table) MarshalJSONFields() ([]byte, error) {
	return json.Marshal(t.Fields)
}

func (t Table) UnmarshalJSONFields(data []byte) error {
	return json.Unmarshal(data, &t.Fields)
}
