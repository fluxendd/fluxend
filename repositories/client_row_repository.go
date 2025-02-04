package repositories

import (
	"fluxton/models"
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

func (r *ClientRowRepository) Create(tableName string, fields models.Row) (uint64, error) {
	if len(fields) == 0 {
		return 0, fmt.Errorf("no fields to insert")
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

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	var insertedID uint64
	err := r.connection.QueryRow(query, values...).Scan(&insertedID)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %v", err)
	}

	return insertedID, nil
}

func (r *ClientRowRepository) GetByID(tableName string, id uint64) (models.Row, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

	resultRow := make(models.Row)
	row := r.connection.QueryRowx(query, id)

	err := row.MapScan(resultRow)
	if err != nil {
		return nil, fmt.Errorf("error getting row: %v", err)
	}

	return resultRow, nil
}
