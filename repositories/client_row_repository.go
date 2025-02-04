package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ClientRowRepository struct {
	connection *sqlx.DB
}

// NewClientRowRepository A repository for dynamically inserting rows into dynamic tables
func NewClientRowRepository(connection *sqlx.DB) (*ClientRowRepository, error) {
	return &ClientRowRepository{connection: connection}, nil
}

func (r *ClientRowRepository) Create(tableName string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to insert")
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	i := 1
	for col, val := range fields {
		columns = append(columns, col)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, val)
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := r.connection.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error inserting row: %v", err)
	}

	return nil
}
