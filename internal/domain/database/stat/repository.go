package stat

type Repository interface {
	GetTotalDatabaseSize() (string, error)
	GetTotalIndexSize() (string, error)
	GetUnusedIndexes() ([]UnusedIndex, error)
	GetSlowQueries() ([]SlowQuery, error)
	GetIndexScansPerTable() ([]IndexScan, error)
	GetSizePerTable() ([]TableSize, error)
	GetRowCountPerTable() ([]TableRowCount, error)
}
