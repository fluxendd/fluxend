package models

import (
	"fluxton/types"
	"time"
)

type DatabaseStat struct {
	Id           int                   `db:"id" json:"id"`
	DatabaseName string                `db:"database_name" json:"databaseName"`
	TotalSize    string                `db:"total_size" json:"totalSize"`
	IndexSize    string                `db:"index_size" json:"indexSize"`
	UnusedIndex  []types.UnusedIndex   `db:"unused_index" json:"unusedIndex"`
	TableCount   []types.TableRowCount `db:"table_count" json:"tableCount"`
	TableSize    []types.TableSize     `db:"table_size" json:"tableSize"`
	CreatedAt    time.Time             `db:"created_at" json:"createdAt"`
}
