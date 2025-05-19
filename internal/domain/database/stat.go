package database

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

type UnusedIndex struct {
	TableName  string `db:"table_name"`
	IndexName  string `db:"index_name"`
	IndexScans int    `db:"index_scans"`
	IndexSize  string `db:"index_size"`
}

type SlowQuery struct {
	Query     string  `db:"query"`
	Calls     int     `db:"calls"`
	TotalTime float64 `db:"total_time"`
	MeanTime  float64 `db:"mean_time"`
}

type TableSize struct {
	TableName string `db:"table_name"`
	TotalSize string `db:"total_size"`
}

type TableRowCount struct {
	TableName         string `db:"table_name"`
	EstimatedRowCount int    `db:"estimated_row_count"`
}

type IndexScan struct {
	TableName  string `db:"table_name"`
	IndexScans int    `db:"index_scans"`
}
