package stat

import (
	"fluxend/internal/domain/stats"
	"time"
)

type Response struct {
	Id           int                   `db:"id" json:"id"`
	DatabaseName string                `db:"database_name" json:"databaseName"`
	TotalSize    string                `db:"total_size" json:"totalSize"`
	IndexSize    string                `db:"index_size" json:"indexSize"`
	UnusedIndex  []stats.UnusedIndex   `db:"unused_index" json:"unusedIndex"`
	TableCount   []stats.TableRowCount `db:"table_count" json:"tableCount"`
	TableSize    []stats.TableSize     `db:"table_size" json:"tableSize"`
	CreatedAt    time.Time             `db:"created_at" json:"createdAt"`
}
