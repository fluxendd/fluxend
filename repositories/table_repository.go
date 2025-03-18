package repositories

import (
	"fluxton/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TableRepository struct {
	connection *sqlx.DB
}

func NewTableRepository(connection *sqlx.DB) (*TableRepository, error) {
	return &TableRepository{connection: connection}, nil
}

func (r *TableRepository) Exists(name string) (bool, error) {
	var count int
	err := r.connection.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TableRepository) Create(name string, columns []models.Column) error {
	query := "CREATE TABLE " + name + " ("

	for _, column := range columns {
		query += column.Name + " " + column.Type

		if column.Primary {
			query += " PRIMARY KEY"
		}

		if column.Unique {
			query += " UNIQUE"
		}

		if column.NotNull {
			query += " NOT NULL"
		}

		if column.Default != "" {
			query += " DEFAULT " + column.Default
		}

		query += ", "
	}

	query = query[:len(query)-2] + ")"

	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) Duplicate(oldName string, newName string) error {
	_, err := r.connection.Exec(fmt.Sprintf("CREATE TABLE %s AS TABLE %s", newName, oldName))
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) List() ([]models.Table, error) {
	var tables []models.Table
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
		return []models.Table{}, err
	}

	return tables, nil
}

func (r *TableRepository) GetByNameInSchema(schema, name string) (models.Table, error) {
	var table models.Table
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

	err := r.connection.Get(&table, query, schema, name)
	if err != nil {
		return models.Table{}, err
	}

	return table, nil
}

func (r *TableRepository) DropIfExists(name string) error {
	_, err := r.connection.Exec("DROP TABLE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) Rename(oldName string, newName string) error {
	_, err := r.connection.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldName, newName))
	if err != nil {
		return err
	}

	return nil
}
