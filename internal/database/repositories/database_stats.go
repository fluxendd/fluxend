package repositories

import (
	"fluxton/internal/domain/database/stat"
	"github.com/jmoiron/sqlx"
)

type DatabaseStatsRepository struct {
	connection *sqlx.DB
}

func NewDatabaseStatsRepository(connection *sqlx.DB) (*DatabaseStatsRepository, error) {
	return &DatabaseStatsRepository{connection: connection}, nil
}

func (r *DatabaseStatsRepository) GetTotalDatabaseSize() (string, error) {
	var totalSize string
	err := r.connection.Get(&totalSize, "SELECT pg_size_pretty(pg_database_size(current_database())) AS database_size;")
	if err != nil {
		return "", err
	}

	return totalSize, nil
}

func (r *DatabaseStatsRepository) GetTotalIndexSize() (string, error) {
	var totalSize string
	err := r.connection.Get(&totalSize, "SELECT pg_size_pretty(sum(pg_relation_size(indexrelid))) AS total_index_size FROM pg_stat_user_indexes;")
	if err != nil {
		return "", err
	}

	return totalSize, nil
}

func (r *DatabaseStatsRepository) GetUnusedIndexes() ([]stat.UnusedIndex, error) {
	var unusedIndexes []stat.UnusedIndex
	err := r.connection.Select(&unusedIndexes, `
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

func (r *DatabaseStatsRepository) GetSlowQueries() ([]stat.SlowQuery, error) {
	var slowQueries []stat.SlowQuery
	err := r.connection.Select(&slowQueries, `
		SELECT query, calls, total_time, mean_time
		FROM pg_stat_statements
		ORDER BY total_time DESC
		LIMIT 5;
	`)
	if err != nil {
		return nil, err
	}

	return slowQueries, nil
}

func (r *DatabaseStatsRepository) GetIndexScansPerTable() ([]stat.IndexScan, error) {
	var indexScans []stat.IndexScan
	err := r.connection.Select(&indexScans, `
		SELECT relname AS table_name, idx_scan AS index_scans
		FROM pg_stat_user_tables
		ORDER BY idx_scan DESC;
	`)
	if err != nil {
		return nil, err
	}

	return indexScans, nil
}

func (r *DatabaseStatsRepository) GetSizePerTable() ([]stat.TableSize, error) {
	var tableSizes []stat.TableSize
	err := r.connection.Select(&tableSizes, `
		SELECT 
			relname AS table_name, 
			pg_size_pretty(pg_total_relation_size(relid)) AS total_size
		FROM pg_catalog.pg_statio_user_tables
		ORDER BY pg_total_relation_size(relid) DESC;
	`)
	if err != nil {
		return nil, err
	}

	return tableSizes, nil
}

func (r *DatabaseStatsRepository) GetRowCountPerTable() ([]stat.TableRowCount, error) {
	var rowCounts []stat.TableRowCount
	err := r.connection.Select(&rowCounts, `
		SELECT 
			relname AS table_name, 
			n_live_tup AS estimated_row_count
		FROM pg_stat_user_tables
		ORDER BY estimated_row_count DESC;
	`)
	if err != nil {
		return nil, err
	}

	return rowCounts, nil
}
