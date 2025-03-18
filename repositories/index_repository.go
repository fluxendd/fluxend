package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type IndexRepository struct {
	connection *sqlx.DB
}

func NewIndexRepository(connection *sqlx.DB) (*IndexRepository, error) {
	return &IndexRepository{connection: connection}, nil
}

func (r *IndexRepository) GetByName(tableName string, indexName string) (string, error) {
	var index string
	err := r.connection.Get(&index, fmt.Sprintf("SELECT indexdef FROM pg_indexes WHERE tablename = '%s' AND indexname = '%s'", tableName, indexName))
	if err != nil {
		return "", err
	}

	return index, nil
}

func (r *IndexRepository) Has(tableName string, indexName string) (bool, error) {
	var count int
	err := r.connection.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM pg_indexes WHERE tablename = '%s' AND indexname = '%s'", tableName, indexName))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *IndexRepository) List(tableName string) ([]string, error) {
	var indexes []string
	err := r.connection.Select(&indexes, fmt.Sprintf("SELECT indexname FROM pg_indexes WHERE tablename = '%s'", tableName))
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

	query := fmt.Sprintf("CREATE %s INDEX %s ON %s (%s)", unique, indexName, tableName, columnsStr)

	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *IndexRepository) DropIfExists(indexName string) (bool, error) {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)

	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}
