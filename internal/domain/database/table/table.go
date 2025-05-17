package table

import (
	"fluxton/internal/domain/shared"
)

type Table struct {
	shared.BaseModel
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Schema        string `db:"schema"`
	EstimatedRows int    `db:"estimated_rows"`
	TotalSize     string `db:"total_size"`
}
