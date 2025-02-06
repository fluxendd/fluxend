package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ClientIndexRepository struct {
	connection *sqlx.DB
}

func NewClientIndexRepository(connection *sqlx.DB) (*ClientColumnRepository, error) {
	return &ClientColumnRepository{connection: connection}, nil
}

func (r *ClientIndexRepository) GetByName(tableName string, indexName string) (string, error) {
	var index string
	err := r.connection.Get(&index, fmt.Sprintf("SELECT indexdef FROM pg_indexes WHERE tablename = '%s' AND indexname = '%s'", tableName, indexName))
	if err != nil {
		return "", err
	}

	return index, nil
}

func (r *ClientIndexRepository) List(tableName string) ([]string, error) {
	var indexes []string
	err := r.connection.Select(&indexes, fmt.Sprintf("SELECT indexname FROM pg_indexes WHERE tablename = '%s'", tableName))
	if err != nil {
		return []string{}, err
	}

	return indexes, nil
}

func (r *ClientIndexRepository) Create(tableName string, indexName string, isUnique bool, columns []string) (bool, error) {
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

func (r *ClientIndexRepository) DropIfExists(indexName string) (bool, error) {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)

	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}
