package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/types"
	"fluxton/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type CoreTableRepository struct {
	db *sqlx.DB
}

func NewCoreTableRepository(injector *do.Injector) (*CoreTableRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &CoreTableRepository{db: db}, nil
}

func (r *CoreTableRepository) ListForProject(paginationParams utils.PaginationParams, projectID uint) ([]models.Table, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	modelSkeleton := models.Table{}

	query := `
		SELECT 
			%s 
		FROM 
			tables
		WHERE 
			project_id = :project_id
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, modelSkeleton.GetColumnsWithAlias(modelSkeleton.GetTableName()))

	params := map[string]interface{}{
		"project_id": projectID,
		"sort":       paginationParams.Sort,
		"limit":      paginationParams.Limit,
		"offset":     offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var table models.Table
		if err := rows.StructScan(&table); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return tables, nil
}

func (r *CoreTableRepository) ListColumns(tableID uint) ([]types.TableColumn, error) {
	query := "SELECT columns FROM tables WHERE id = $1"

	var columnsJSON models.JSONColumns
	row := r.db.QueryRow(query, tableID)

	err := row.Scan(&columnsJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("table.error.notFound")
		}
		return nil, fmt.Errorf("could not fetch row: %v", err)
	}

	return columnsJSON, nil
}

func (r *CoreTableRepository) GetByID(id uint) (models.Table, error) {
	query := "SELECT %s FROM tables WHERE id = $1"
	query = fmt.Sprintf(query, models.Table{}.GetColumns())

	var table models.Table
	err := r.db.Get(&table, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Table{}, errs.NewNotFoundError("table.error.notFound")
		}
		return models.Table{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return table, nil
}

func (r *CoreTableRepository) GetByName(name string) (models.Table, error) {
	query := "SELECT %s FROM tables WHERE name = $1"
	query = fmt.Sprintf(query, models.Table{}.GetColumns())

	var table models.Table
	err := r.db.Get(&table, query, name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Table{}, errs.NewNotFoundError("table.error.notFound")
		}
		return models.Table{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return table, nil
}

func (r *CoreTableRepository) ExistsByID(id uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tables WHERE id = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *CoreTableRepository) ExistsByNameForProject(name string, tableID uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tables WHERE name = $1 AND project_Id = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, tableID)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *CoreTableRepository) HasColumn(column string, tableID uint) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tables
			WHERE id = $1
			AND EXISTS (
				SELECT 1
				FROM jsonb_array_elements(columns) AS col
				WHERE col->>'Name' = $2
			)
		) AS column_exists
	`

	var columnExists bool
	err := r.db.QueryRow(query, tableID, column).Scan(&columnExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.NewNotFoundError("table.error.notFound")
		}
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return columnExists, nil
}

func (r *CoreTableRepository) Create(table *models.Table) (*models.Table, error) {
	columnsJSON, err := table.MarshalJSONColumns()
	if err != nil {
		return nil, fmt.Errorf("could not marshal columns: %v", err)
	}

	query := "INSERT INTO tables (name, project_id, created_by, updated_by, columns) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	queryErr := r.db.QueryRow(query, table.Name, table.ProjectID, table.CreatedBy, table.UpdatedBy, columnsJSON).Scan(&table.ID)
	if queryErr != nil {
		return nil, fmt.Errorf("could not create table: %v", queryErr)
	}

	return table, nil
}

func (r *CoreTableRepository) Update(table *models.Table) (*models.Table, error) {
	columnsJSON, err := table.MarshalJSONColumns()
	if err != nil {
		return nil, fmt.Errorf("could not marshal columns: %v", err)
	}

	query := `
		UPDATE tables 
		SET name = $1, columns = $2, updated_at = $3, updated_by = $4
		WHERE id = $5
		RETURNING id
	`

	queryErr := r.db.QueryRow(query, table.Name, columnsJSON, time.Now(), table.UpdatedBy, table.ID).Scan(&table.ID)
	if queryErr != nil {
		return nil, fmt.Errorf("could not update table: %v", queryErr)
	}

	return table, nil
}

func (r *CoreTableRepository) Rename(tableID uint, name string, authenticatedUserID uint) (models.Table, error) {
	query := `
		UPDATE tables 
		SET name = $1, updated_at = $2, updated_by = $3
		WHERE id = $4`

	queryErr := r.db.QueryRow(query, name, time.Now(), authenticatedUserID, tableID)
	if queryErr != nil {
		return models.Table{}, fmt.Errorf("could not update table: %v", queryErr)
	}

	return r.GetByID(tableID)
}

func (r *CoreTableRepository) Delete(tableID uint) (bool, error) {
	query := "DELETE FROM tables WHERE id = $1"
	res, err := r.db.Exec(query, tableID)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
