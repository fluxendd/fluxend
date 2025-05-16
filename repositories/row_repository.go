package repositories

import (
	"fluxton/models"
	"github.com/jmoiron/sqlx"
)

type RowRepository struct {
	connection *sqlx.DB
}

func NewRowRepository(connection *sqlx.DB) (*RowRepository, error) {
	return &RowRepository{connection: connection}, nil
}

func (r *RowRepository) CreateMany(tableName string, columns []models.Column, values [][]string) error {
	query := "INSERT INTO " + tableName + " ("

	for i, column := range columns {
		query += column.Name
		if i < len(columns)-1 {
			query += ", "
		}
	}

	query += ") VALUES "

	for i, valueSet := range values {
		query += "("
		for j, value := range valueSet {
			query += "'" + value + "'"
			if j < len(valueSet)-1 {
				query += ", "
			}
		}
		query += ")"
		if i < len(values)-1 {
			query += ", "
		}
	}

	_, err := r.connection.Exec(query)

	return err
}
