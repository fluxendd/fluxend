package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/domain/database"
	flxErrs "fluxton/pkg/errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type TableRepository struct {
	connection       *sqlx.DB
	columnRepository *ColumnRepository
}

func NewTableRepository(connection *sqlx.DB) (database.Repository, error) {
	columnRepository, err := NewColumnRepository(connection)
	if err != nil {
		return nil, err
	}

	return &TableRepository{
		connection:       connection,
		columnRepository: columnRepository,
	}, nil
}

func (r *TableRepository) Exists(name string) (bool, error) {
	var count int
	err := r.connection.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TableRepository) Create(name string, columns []database.Column) error {
	var defs []string
	var foreignConstraints []string

	for _, currentColumn := range columns {
		defs = append(defs, r.columnRepository.BuildColumnDefinition(currentColumn))

		if fkQuery, ok := r.columnRepository.BuildForeignKeyConstraint(name, currentColumn); ok {
			foreignConstraints = append(foreignConstraints, fkQuery)
		}
	}

	createQuery := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", pq.QuoteIdentifier(name), strings.Join(defs, ",\n"))

	if _, err := r.connection.Exec(createQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	for _, fk := range foreignConstraints {
		if _, err := r.connection.Exec(fk); err != nil {
			return fmt.Errorf("failed to add foreign key constraint: %w", err)
		}
	}

	return nil
}

func (r *TableRepository) Duplicate(existingTable string, newTable string) error {
	_, err := r.connection.Exec(fmt.Sprintf("CREATE TABLE %s AS TABLE %s", pq.QuoteIdentifier(newTable), pq.QuoteIdentifier(existingTable)))
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) List() ([]database.Table, error) {
	var tables []database.Table
	query := `
		SELECT
			c.oid AS id,
			c.relname AS name,
			n.nspname AS schema,
			c.reltuples AS estimated_rows,  -- Approximate row count
			pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size -- Table size (including indexes)
		FROM pg_class c
				 JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = 'public'  -- Only list tables in the "public" schema
		  AND c.relkind = 'r'  -- 'r' means regular table (excludes views, indexes, etc.)
		ORDER BY c.relname;
	`
	err := r.connection.Select(&tables, query)
	if err != nil {
		return []database.Table{}, err
	}

	return tables, nil
}

func (r *TableRepository) GetByNameInSchema(schema, name string) (database.Table, error) {
	var fetchedTable database.Table
	query := `
		SELECT
			c.oid AS id,
			c.relname AS name,
			n.nspname AS schema,
			c.reltuples AS estimated_rows,  -- Approximate row count
			pg_size_pretty(pg_total_relation_size(c.oid)) AS total_size -- Table size (including indexes)
		FROM pg_class c
		JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE n.nspname = $1  -- Filter by schema
		  AND c.relname = $2  -- Filter by table name
		LIMIT 1;
	`

	err := r.connection.Get(&fetchedTable, query, schema, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Table{}, flxErrs.NewNotFoundError("table.error.notFound")
		}

		return database.Table{}, err
	}

	return fetchedTable, nil
}

func (r *TableRepository) DropIfExists(name string) error {
	_, err := r.connection.Exec("DROP TABLE IF EXISTS " + pq.QuoteIdentifier(name))
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) Rename(oldName string, newName string) error {
	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", pq.QuoteIdentifier(oldName), pq.QuoteIdentifier(newName))
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
