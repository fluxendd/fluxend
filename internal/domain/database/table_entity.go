package database

import (
	"fluxend/internal/domain/shared"
)

type Table struct {
	shared.BaseEntity
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Schema        string `db:"schema"`
	EstimatedRows int    `db:"estimated_rows"`
	TotalSize     string `db:"total_size"`
}
