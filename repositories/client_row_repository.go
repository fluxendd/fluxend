package repositories

import (
	"fluxton/models"
	"fluxton/utils"
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

func (r *CoreTableRepository) List(tableName string, paginationParams utils.PaginationParams) ([]models.Row, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT * FROM %s
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, tableName)

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var results []models.Row
	for rows.Next() {
		var result models.Row
		if err := rows.StructScan(&result); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return results, nil
}

func (r *ClientRowRepository) GetByID(tableName string, rowID uint64) (models.Row, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

	resultRow := make(models.Row)
	row := r.connection.QueryRowx(query, rowID)

	err := row.MapScan(resultRow)
	if err != nil {
		return nil, fmt.Errorf("error getting row: %v", err)
	}

	return resultRow, nil
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

func (r *ClientRowRepository) Update(tableName string, rowID uint64, fields models.Row) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	var columns []string
	var values []interface{}

	i := 1
	for col, val := range fields {
		columns = append(columns, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		tableName,
		strings.Join(columns, ", "),
		i,
	)
	values = append(values, rowID)

	_, err := r.connection.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error updating row: %v", err)
	}

	return nil
}

func (r *ClientRowRepository) Delete(tableName string, rowID uint64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	_, err := r.connection.Exec(query, rowID)
	if err != nil {
		return fmt.Errorf("error deleting row: %v", err)
	}

	return nil
}
