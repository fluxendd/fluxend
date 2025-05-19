package repositories

import (
	"fluxton/internal/domain/database"
	"github.com/jmoiron/sqlx"
)

type RowRepository struct {
	connection *sqlx.DB
}

func NewRowRepository(connection *sqlx.DB) (database.Repository, error) {
	return &RowRepository{connection: connection}, nil
}

func (r *RowRepository) CreateMany(tableName string, columns []database.Column, values [][]string) error {
	query := "INSERT INTO " + tableName + " ("

	for i, currentColumn := range columns {
		query += currentColumn.Name
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
