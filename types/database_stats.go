package types

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
