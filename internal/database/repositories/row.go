package repositories

import (
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/shared"
	"github.com/samber/do"
)

type RowRepository struct {
	db shared.DB
}

func NewRowRepository(injector *do.Injector) (database.RowRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &RowRepository{db: db}, nil
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

	_, err := r.db.ExecWithRowsAffected(query)
	return err
}
