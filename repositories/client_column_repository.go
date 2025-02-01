package repositories

import (
	"fluxton/types"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ClientColumnRepository struct {
	connection *sqlx.DB
}

func NewClientColumnRepository(connection *sqlx.DB) (*ClientColumnRepository, error) {
	return &ClientColumnRepository{connection: connection}, nil
}

func (r *ClientColumnRepository) List(name string) ([]string, error) {
	var columns []string
	err := r.connection.Select(&columns, fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", name))
	if err != nil {
		return []string{}, err
	}

	return columns, nil
}

func (r *ClientColumnRepository) Create(name string, field types.TableColumn) (bool, error) {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", name, field.Name, field.Type)
	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ClientColumnRepository) Update(name string, field types.TableColumn) (bool, error) {
	query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", name, field.Name, field.Type)
	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ClientColumnRepository) Rename(name, oldColumnName, newColumnName string) (bool, error) {
	query := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", name, oldColumnName, newColumnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ClientColumnRepository) Drop(name, columnName string) (bool, error) {
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", name, columnName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return false, err
	}

	return true, nil
}
