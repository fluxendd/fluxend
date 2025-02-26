package models

import (
	"encoding/json"
	"fluxton/types"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type JSONColumns []types.TableColumn // important for reading from db

type Table struct {
	Uuid        uuid.UUID   `db:"uuid"`
	ProjectUuid uuid.UUID   `db:"project_uuid"`
	CreatedBy   uuid.UUID   `db:"created_by"`
	UpdatedBy   uuid.UUID   `db:"updated_by"`
	Name        string      `db:"name"`
	Columns     JSONColumns `db:"columns" json:"columns"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

func (t Table) GetTableName() string {
	return "fluxton.tables"
}

func (t Table) MarshalJSONColumns() ([]byte, error) {
	return json.Marshal(t.Columns)
}

func (t Table) UnmarshalJSONColumns(data []byte) error {
	return json.Unmarshal(data, &t.Columns)
}

func (j *JSONColumns) Scan(value interface{}) error {
	byteData, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert database value to []byte")
	}
	return json.Unmarshal(byteData, j)
}
