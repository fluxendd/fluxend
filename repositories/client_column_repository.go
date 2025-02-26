package repositories

import (
	"fluxton/types"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ClientColumnRepository struct {
	connection *sqlx.DB
}

func NewClientColumnRepository(connection *sqlx.DB) (*ClientColumnRepository, error) {
	return &ClientColumnRepository{connection: connection}, nil
}

func (r *ClientColumnRepository) List(tableName string) ([]string, error) {
	var columns []string
	err := r.connection.Select(&columns, fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", tableName))
	if err != nil {
		return []string{}, err
	}

	return columns, nil
}

func (r *ClientColumnRepository) Has(tableName, columnName string) (bool, error) {
	var count int
	err := r.connection.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = '%s' AND column_name = '%s'", tableName, columnName))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ClientColumnRepository) HasAny(tableName string, columnNames []string) (bool, error) {
	var count int
	err := r.connection.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = '%s' AND column_name IN ('%s')", tableName, columnNames))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ClientColumnRepository) HasAll(tableName string, columns []types.TableColumn) (bool, error) {
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

func (r *ClientColumnRepository) Create(tableName string, field types.TableColumn) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, field.Name, field.Type)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientColumnRepository) CreateMany(tableName string, fields []types.TableColumn) error {
	for _, field := range fields {
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, field.Name, field.Type)
		_, err := r.connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ClientColumnRepository) Alter(tableName string, columns []types.TableColumn) error {
	for _, column := range columns {
		query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, column.Name, column.Type)
		_, err := r.connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ClientColumnRepository) AlterMany(tableName string, fields []types.TableColumn) error {
	for _, field := range fields {
		query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, field.Name, field.Type)
		_, err := r.connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ClientColumnRepository) Rename(tableName, oldColumnName, newColumnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tableName, oldColumnName, newColumnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientColumnRepository) Drop(tableName, columnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientColumnRepository) mapColumnsToNames(columns []types.TableColumn) []string {
	columnNames := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = column.Name
	}

	return columnNames
}
