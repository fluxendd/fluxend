package repositories

import (
	"fluxend/internal/domain/shared"
	"fmt"
	"github.com/lib/pq"
	"github.com/samber/do"
	"strings"
)

type IndexRepository struct {
	db shared.DB
}

func NewIndexRepository(injector *do.Injector) (*IndexRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &IndexRepository{db: db}, nil
}

func (r *IndexRepository) GetByName(tableName, indexName string) (string, error) {
	var index string
	query := `
       SELECT indexdef
       FROM pg_indexes
       WHERE tablename = $1 AND indexname = $2
    `
	return index, r.db.Get(&index, query, tableName, indexName)
}

func (r *IndexRepository) Has(tableName, indexName string) (bool, error) {
	return r.db.Exists("pg_indexes", "tablename = $1 AND indexname = $2", tableName, indexName)
}

func (r *IndexRepository) List(tableName string) ([]string, error) {
	var indexes []string
	query := `
       SELECT indexname
       FROM pg_indexes
       WHERE tablename = $1
    `
	return indexes, r.db.Select(&indexes, query, tableName)
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

	_, err := r.db.ExecWithRowsAffected(query)
	return err == nil, err
}

func (r *IndexRepository) DropIfExists(indexName string) (bool, error) {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", pq.QuoteIdentifier(indexName))

	_, err := r.db.ExecWithRowsAffected(query)
	return err == nil, err
}
