package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type IndexRepository struct {
	connection *sqlx.DB
}

func NewIndexRepository(connection *sqlx.DB) (*IndexRepository, error) {
	return &IndexRepository{connection: connection}, nil
}

func (r *IndexRepository) GetByName(tableName, indexName string) (string, error) {
	var index string
	query := `
		SELECT indexdef
		FROM pg_indexes
		WHERE tablename = $1 AND indexname = $2
	`
	err := r.connection.Get(&index, query, tableName, indexName)
	if err != nil {
		return "", err
	}

	return index, nil
}

func (r *IndexRepository) Has(tableName, indexName string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE tablename = $1 AND indexname = $2
	`
	err := r.connection.Get(&count, query, tableName, indexName)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *IndexRepository) List(tableName string) ([]string, error) {
	var indexes []string
	query := `
		SELECT indexname
		FROM pg_indexes
		WHERE tablename = $1
	`
	err := r.connection.Select(&indexes, query, tableName)
	if err != nil {
		return []string{}, err
	}

	return indexes, nil
}

func (r *IndexRepository) Create(tableName string, indexName string, columns []string, isUnique bool) (bool, error) {
	unique := ""
	if isUnique {
		unique = "UNIQUE"
	}

	columnsStr := strings.Join(columns, ", ")

	query := fmt.Sprintf(
		"CREATE %s INDEX %s ON %s (%s)",
		unique,
		pq.QuoteIdentifier(indexName),
		pq.QuoteIdentifier(tableName),
		columnsStr,
	)

	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *IndexRepository) DropIfExists(indexName string) (bool, error) {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", pq.QuoteIdentifier(indexName))

	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}
