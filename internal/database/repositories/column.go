package repositories

import (
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/shared"
	"fmt"
	"github.com/lib/pq"
	"github.com/samber/do"
	"strings"
)

type ColumnRepository struct {
	db shared.DB
}

func NewColumnRepository(injector *do.Injector) (database.ColumnRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &ColumnRepository{db: db}, nil
}

func (r *ColumnRepository) List(tableName string) ([]database.Column, error) {
	var columns []database.Column
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
	return columns, r.db.Select(&columns, query, tableName)
}

func (r *ColumnRepository) Has(tableName, columnName string) (bool, error) {
	return r.db.Exists("information_schema.columns", "table_name = $1 AND column_name = $2", tableName, columnName)
}

func (r *ColumnRepository) HasAny(tableName string, columns []database.Column) (bool, error) {
	var count int
	columnNames := r.mapColumnsToNames(columns)
	query := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = $1
		AND column_name = ANY($2)
	`

	err := r.db.Get(&count, query, tableName, pq.Array(columnNames))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ColumnRepository) HasAll(tableName string, columns []database.Column) (bool, error) {
	var count int
	columnNames := r.mapColumnsToNames(columns)
	query := `
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = $1 
		AND column_name = ANY($2)
	`

	err := r.db.Get(&count, query, tableName, pq.Array(columnNames))
	if err != nil {
		return false, err
	}

	return count == len(columns), nil
}

func (r *ColumnRepository) CreateOne(tableName string, column database.Column) error {
	return r.db.WithTransaction(func(tx shared.Tx) error {
		def := r.BuildColumnDefinition(column)
		addColumnQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", tableName, def)

		if _, err := tx.Exec(addColumnQuery); err != nil {
			return fmt.Errorf("failed to add column: %w", err)
		}

		if fkQuery, ok := r.BuildForeignKeyConstraint(tableName, column); ok {
			if _, err := tx.Exec(fkQuery); err != nil {
				return fmt.Errorf("failed to add foreign key constraint: %w", err)
			}
		}

		return nil
	})
}

func (r *ColumnRepository) CreateMany(tableName string, fields []database.Column) error {
	for _, field := range fields {
		if err := r.CreateOne(tableName, field); err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) AlterOne(tableName string, columns []database.Column) error {
	return r.db.WithTransaction(func(tx shared.Tx) error {
		for _, column := range columns {
			query := fmt.Sprintf(
				"ALTER TABLE %s ALTER COLUMN %s TYPE %s",
				tableName,
				column.Name,
				column.Type,
			)

			if _, err := tx.Exec(query); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ColumnRepository) AlterMany(tableName string, fields []database.Column) error {
	for _, field := range fields {
		if err := r.AlterOne(tableName, []database.Column{field}); err != nil {
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

	return r.db.ExecWithErr(query)
}

func (r *ColumnRepository) Drop(tableName, columnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, pq.QuoteIdentifier(columnName))
	_, err := r.db.ExecWithRowsAffected(query)
	return err
}

func (r *ColumnRepository) DropMany(tableName string, columns []database.Column) error {
	if len(columns) == 0 {
		return fmt.Errorf("no columns specified")
	}

	var drops []string
	for _, column := range columns {
		drops = append(drops, fmt.Sprintf("DROP COLUMN %s", pq.QuoteIdentifier(column.Name)))
	}

	query := fmt.Sprintf("ALTER TABLE %s %s",
		pq.QuoteIdentifier(tableName),
		strings.Join(drops, ", "))

	_, err := r.db.ExecWithRowsAffected(query)

	return err
}

func (r *ColumnRepository) mapColumnsToNames(columns []database.Column) []string {
	columnNames := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = column.Name
	}

	return columnNames
}

func (r *ColumnRepository) BuildColumnDefinition(column database.Column) string {
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

func (r *ColumnRepository) BuildForeignKeyConstraint(tableName string, column database.Column) (string, bool) {
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
