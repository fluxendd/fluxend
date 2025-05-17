package repositories

import (
	"fluxton/internal/domain/database/column"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ColumnRepository struct {
	connection *sqlx.DB
}

func NewColumnRepository(connection *sqlx.DB) (*ColumnRepository, error) {
	return &ColumnRepository{connection: connection}, nil
}

func (r *ColumnRepository) List(tableName string) ([]column.Column, error) {
	var columns []column.Column
	query := `
		SELECT 
			a.attname AS name,
			a.attnum AS position,
			a.attnotnull AS not_null,
			COALESCE(pg_catalog.format_type(a.atttypid, a.atttypmod), '') AS type,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
			COALESCE(ct.contype = 'p', false) AS primary,
			COALESCE(ct.contype = 'u', false) AS unique,
			COALESCE(ct.contype = 'f', false) AS foreign,
			ref_table.relname AS reference_table,
			ref_col.attname AS reference_column
		FROM pg_attribute a
		LEFT JOIN pg_attrdef ad 
			ON a.attrelid = ad.adrelid AND a.attnum = ad.adnum
		LEFT JOIN pg_constraint ct 
			ON ct.conrelid = a.attrelid AND a.attnum = ANY(ct.conkey)
		LEFT JOIN pg_class ref_table 
			ON ref_table.oid = ct.confrelid
		LEFT JOIN pg_attribute ref_col 
			ON ref_col.attrelid = ct.confrelid AND ref_col.attnum = ct.confkey[1]
		WHERE a.attrelid = $1::regclass
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		ORDER BY a.attnum;

	`
	err := r.connection.Select(&columns, query, tableName)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *ColumnRepository) Has(tableName, columnName string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2
	`
	err := r.connection.Get(&count, query, tableName, columnName)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ColumnRepository) HasAny(tableName string, columns []column.Column) (bool, error) {
	var count int
	columnNames := r.mapColumnsToNames(columns)
	query := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = $1
		AND column_name = ANY($2)
	`

	err := r.connection.Get(&count, query, tableName, pq.Array(columnNames))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ColumnRepository) HasAll(tableName string, columns []column.Column) (bool, error) {
	var count int
	columnNames := r.mapColumnsToNames(columns)
	query := `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = $1 
		AND column_name = ANY($2)
	`

	err := r.connection.Get(&count, query, tableName, pq.Array(columnNames))
	if err != nil {
		return false, err
	}

	return count == len(columns), nil
}

func (r *ColumnRepository) CreateOne(tableName string, column column.Column) error {
	def := r.BuildColumnDefinition(column)
	addColumnQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", tableName, def)

	if _, err := r.connection.Exec(addColumnQuery); err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}

	if fkQuery, ok := r.BuildForeignKeyConstraint(tableName, column); ok {
		if _, err := r.connection.Exec(fkQuery); err != nil {
			return fmt.Errorf("failed to add foreign key constraint: %w", err)
		}
	}

	return nil
}

func (r *ColumnRepository) CreateMany(tableName string, fields []column.Column) error {
	for _, field := range fields {
		err := r.CreateOne(tableName, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) AlterOne(tableName string, columns []column.Column) error {
	for _, column := range columns {
		query := fmt.Sprintf(
			"ALTER TABLE %s ALTER COLUMN %s TYPE %s",
			tableName,
			pq.QuoteIdentifier(column.Name),
			pq.QuoteIdentifier(column.Type),
		)
		_, err := r.connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) AlterMany(tableName string, fields []column.Column) error {
	for _, field := range fields {
		err := r.AlterOne(tableName, []column.Column{field})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) Rename(tableName, oldColumnName, newColumnName string) error {
	query := fmt.Sprintf(
		"ALTER TABLE %s RENAME COLUMN %s TO %s",
		tableName,
		pq.QuoteIdentifier(oldColumnName),
		pq.QuoteIdentifier(newColumnName),
	)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ColumnRepository) Drop(tableName, columnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, pq.QuoteIdentifier(columnName))
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ColumnRepository) mapColumnsToNames(columns []column.Column) []string {
	columnNames := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = column.Name
	}

	return columnNames
}

func (r *ColumnRepository) BuildColumnDefinition(column column.Column) string {
	def := fmt.Sprintf("%s %s", pq.QuoteIdentifier(column.Name), column.Type)

	if column.Primary {
		def += " PRIMARY KEY"
	}

	if column.Unique {
		def += " UNIQUE"
	}

	if column.NotNull {
		def += " NOT NULL"
	}

	if column.Default != "" {
		def += fmt.Sprintf(" DEFAULT %s", column.Default)
	}

	return def
}

func (r *ColumnRepository) BuildForeignKeyConstraint(tableName string, column column.Column) (string, bool) {
	if !column.Foreign || !column.ReferenceTable.Valid || !column.ReferenceColumn.Valid {
		return "", false
	}

	return fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT fk_%s FOREIGN KEY (%s) REFERENCES %s(%s);",
		tableName,
		pq.QuoteIdentifier(column.Name),
		pq.QuoteIdentifier(column.Name),
		column.ReferenceTable.String,
		column.ReferenceColumn.String,
	), true
}
