package repositories

import (
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/shared"
	"fluxend/pkg"
	"fmt"
	"github.com/lib/pq"
	"github.com/samber/do"
	"strings"
)

type TableRepository struct {
	db               shared.DB
	columnRepository database.ColumnRepository
}

func NewTableRepository(injector *do.Injector) (database.TableRepository, error) {
	db := do.MustInvoke[shared.DB](injector)

	columnRepository, err := NewColumnRepository(injector)
	if err != nil {
		return nil, err
	}

	return &TableRepository{
		db:               db,
		columnRepository: columnRepository,
	}, nil
}

func (r *TableRepository) Exists(name string) (bool, error) {
	return r.db.Exists("information_schema.tables", "table_schema = 'public' AND table_name = $1", name)
}

func (r *TableRepository) Create(name string, columns []database.Column) error {
	return r.db.WithTransaction(func(tx shared.Tx) error {
		var defs []string
		var foreignConstraints []string

		for _, currentColumn := range columns {
			defs = append(defs, r.columnRepository.BuildColumnDefinition(currentColumn))

			if fkQuery, ok := r.columnRepository.BuildForeignKeyConstraint(name, currentColumn); ok {
				foreignConstraints = append(foreignConstraints, fkQuery)
			}
		}

		createQuery := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", pq.QuoteIdentifier(name), strings.Join(defs, ",\n"))

		if _, err := tx.Exec(createQuery); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		for _, fk := range foreignConstraints {
			if _, err := tx.Exec(fk); err != nil {
				return fmt.Errorf("failed to add foreign key constraint: %w", err)
			}
		}

		return nil
	})
}

func (r *TableRepository) Duplicate(existingTable string, newTable string) error {
	query := fmt.Sprintf("CREATE TABLE %s AS TABLE %s", pq.QuoteIdentifier(newTable), pq.QuoteIdentifier(existingTable))

	return r.db.ExecWithErr(query)
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
	return tables, r.db.Select(&tables, query)
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

	return fetchedTable, r.db.GetWithNotFound(&fetchedTable, "table.error.notFound", query, schema, name)
}

func (r *TableRepository) DropIfExists(name string) error {
	schema, name := pkg.ParseTableName(name)

	return r.db.ExecWithErr(fmt.Sprintf(`DROP TABLE IF EXISTS %s."%s"`, schema, name))
}

func (r *TableRepository) Rename(oldName string, newName string) error {
	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", pq.QuoteIdentifier(oldName), pq.QuoteIdentifier(newName))

	return r.db.ExecWithErr(query)
}
