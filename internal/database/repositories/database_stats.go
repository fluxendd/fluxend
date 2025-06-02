package repositories

import (
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/stats"
	"github.com/samber/do"
)

type DatabaseStatsRepository struct {
	db shared.DB
}

func NewDatabaseStatsRepository(injector *do.Injector) (stats.StatRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &DatabaseStatsRepository{db: db}, nil
}

func (r *DatabaseStatsRepository) GetTotalDatabaseSize() (string, error) {
	var totalSize string
	return totalSize, r.db.Get(&totalSize, "SELECT pg_size_pretty(pg_database_size(current_database())) AS database_size;")
}

func (r *DatabaseStatsRepository) GetTotalIndexSize() (string, error) {
	var totalSize string
	return totalSize, r.db.Get(&totalSize, "SELECT pg_size_pretty(sum(pg_relation_size(indexrelid))) AS total_index_size FROM pg_stat_user_indexes;")
}

func (r *DatabaseStatsRepository) GetUnusedIndexes() ([]stats.UnusedIndex, error) {
	var unusedIndexes []stats.UnusedIndex
	err := r.db.Select(&unusedIndexes, `
		SELECT 
			relname AS table_name, 
			indexrelname AS index_name, 
			idx_scan AS index_scans,
			pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
		FROM pg_stat_user_indexes
		WHERE idx_scan < 50
		ORDER BY idx_scan;
	`)
	if err != nil {
		return nil, err
	}

	return unusedIndexes, nil
}

func (r *DatabaseStatsRepository) GetSlowQueries() ([]stats.SlowQuery, error) {
	var slowQueries []stats.SlowQuery
	query := `
       SELECT query, calls, total_time, mean_time
       FROM pg_stat_statements
       ORDER BY total_time DESC
       LIMIT 5;
    `
	return slowQueries, r.db.Select(&slowQueries, query)
}

func (r *DatabaseStatsRepository) GetIndexScansPerTable() ([]stats.IndexScan, error) {
	var indexScans []stats.IndexScan
	query := `
       SELECT relname AS table_name, idx_scan AS index_scans
       FROM pg_stat_user_tables
       ORDER BY idx_scan DESC;
    `
	return indexScans, r.db.Select(&indexScans, query)
}

func (r *DatabaseStatsRepository) GetSizePerTable() ([]stats.TableSize, error) {
	var tableSizes []stats.TableSize
	query := `
       SELECT 
          relname AS table_name, 
          pg_size_pretty(pg_total_relation_size(relid)) AS total_size
       FROM pg_catalog.pg_statio_user_tables
       ORDER BY pg_total_relation_size(relid) DESC;
    `
	return tableSizes, r.db.Select(&tableSizes, query)
}

func (r *DatabaseStatsRepository) GetRowCountPerTable() ([]stats.TableRowCount, error) {
	var rowCounts []stats.TableRowCount
	query := `
       SELECT 
          relname AS table_name, 
          n_live_tup AS estimated_row_count
       FROM pg_stat_user_tables
       ORDER BY estimated_row_count DESC;
    `
	return rowCounts, r.db.Select(&rowCounts, query)
}
