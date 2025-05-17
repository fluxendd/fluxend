package stat

import (
	"time"
)

type DatabaseStat struct {
	Id           int             `db:"id" json:"id"`
	DatabaseName string          `db:"database_name" json:"databaseName"`
	TotalSize    string          `db:"total_size" json:"totalSize"`
	IndexSize    string          `db:"index_size" json:"indexSize"`
	UnusedIndex  []UnusedIndex   `db:"unused_index" json:"unusedIndex"`
	TableCount   []TableRowCount `db:"table_count" json:"tableCount"`
	TableSize    []TableSize     `db:"table_size" json:"tableSize"`
	CreatedAt    time.Time       `db:"created_at" json:"createdAt"`
}
