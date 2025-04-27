package repositories

import (
	"fluxton/models"
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

func (r *ColumnRepository) List(tableName string) ([]models.Column, error) {
	var columns []models.Column
	query := `
		SELECT 
			a.attname AS column_name,
			a.attnum AS column_position,
			a.attnotnull AS not_null,
			COALESCE(pg_catalog.format_type(a.atttypid, a.atttypmod), '') AS data_type,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
			CASE 
				WHEN ct.contype = 'p' THEN 'PRIMARY'
				WHEN ct.contype = 'u' THEN 'UNIQUE'
			    WHEN ct.contype = 'f' THEN 'FOREIGN'
				ELSE ''
			END AS constraint_type
		FROM pg_attribute a
		LEFT JOIN pg_attrdef ad ON a.attrelid = ad.adrelid AND a.attnum = ad.adnum
		LEFT JOIN pg_constraint ct ON ct.conrelid = a.attrelid AND a.attnum = ANY(ct.conkey)
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

func (r *ColumnRepository) HasAny(tableName string, columns []models.Column) (bool, error) {
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

func (r *ColumnRepository) HasAll(tableName string, columns []models.Column) (bool, error) {
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

func (r *ColumnRepository) CreateOne(tableName string, column models.Column) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, column.Name, column.Type)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ColumnRepository) CreateMany(tableName string, fields []models.Column) error {
	for _, field := range fields {
		err := r.CreateOne(tableName, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) AlterOne(tableName string, columns []models.Column) error {
	for _, column := range columns {
		query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, column.Name, column.Type)
		_, err := r.connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) AlterMany(tableName string, fields []models.Column) error {
	for _, field := range fields {
		err := r.AlterOne(tableName, []models.Column{field})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ColumnRepository) Rename(tableName, oldColumnName, newColumnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tableName, oldColumnName, newColumnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ColumnRepository) Drop(tableName, columnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ColumnRepository) mapColumnsToNames(columns []models.Column) []string {
	columnNames := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = column.Name
	}

	return columnNames
}
