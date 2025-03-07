package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/types"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
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

func (r *CoreTableRepository) ListForProject(paginationParams utils.PaginationParams, projectUUID uuid.UUID) ([]models.Table, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT 
			%s 
		FROM 
			fluxton.tables
		WHERE 
			project_uuid = :project_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumns[models.Table]())

	params := map[string]interface{}{
		"project_uuid": projectUUID,
		"sort":         paginationParams.Sort,
		"limit":        paginationParams.Limit,
		"offset":       offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var table models.Table
		if err := rows.StructScan(&table); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return tables, nil
}

func (r *CoreTableRepository) ListColumns(tableUUID uuid.UUID) ([]types.TableColumn, error) {
	query := "SELECT columns FROM fluxton.tables WHERE uuid = $1"

	var columnsJSON models.JSONColumns
	row := r.db.QueryRow(query, tableUUID)

	err := row.Scan(&columnsJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("table.error.notFound")
		}
		return nil, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return columnsJSON, nil
}

func (r *CoreTableRepository) GetByID(tableUUID uuid.UUID) (models.Table, error) {
	query := "SELECT %s FROM fluxton.tables WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Table]())

	var table models.Table
	err := r.db.Get(&table, query, tableUUID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Table{}, errs.NewNotFoundError("table.error.notFound")
		}
		return models.Table{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return table, nil
}

func (r *CoreTableRepository) GetByName(name string) (models.Table, error) {
	query := "SELECT %s FROM fluxton.tables WHERE name = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Table]())

	var table models.Table
	err := r.db.Get(&table, query, name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Table{}, errs.NewNotFoundError("table.error.notFound")
		}

		return models.Table{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return table, nil
}

func (r *CoreTableRepository) ExistsByID(tableUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.tables WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, tableUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *CoreTableRepository) ExistsByNameForProject(name string, tableUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.tables WHERE name = $1 AND project_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, tableUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *CoreTableRepository) HasColumn(column string, tableUUID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM fluxton.tables
			WHERE uuid = $1
			AND EXISTS (
				SELECT 1
				FROM jsonb_array_elements(columns) AS col
				WHERE col->>'Name' = $2
			)
		) AS column_exists
	`

	var columnExists bool
	err := r.db.QueryRow(query, tableUUID, column).Scan(&columnExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.NewNotFoundError("table.error.notFound")
		}

		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return columnExists, nil
}

func (r *CoreTableRepository) Create(table *models.Table) (*models.Table, error) {
	columnsJSON, err := table.MarshalJSONColumns()
	if err != nil {
		return nil, utils.FormatError(err, "jsonMarshal", utils.GetMethodName())
	}

	query := "INSERT INTO fluxton.tables (name, project_uuid, created_by, updated_by, columns) VALUES ($1, $2, $3, $4, $5) RETURNING uuid"
	queryErr := r.db.QueryRow(query, table.Name, table.ProjectUuid, table.CreatedBy, table.UpdatedBy, columnsJSON).Scan(&table.Uuid)
	if queryErr != nil {
		return nil, utils.FormatError(queryErr, "insert", utils.GetMethodName())
	}

	return table, nil
}

func (r *CoreTableRepository) Update(table *models.Table) (*models.Table, error) {
	columnsJSON, err := table.MarshalJSONColumns()
	if err != nil {
		return nil, fmt.Errorf("could not marshal columns: %v", err)
	}

	query := `
		UPDATE fluxton.tables 
		SET name = $1, columns = $2, updated_at = $3, updated_by = $4
		WHERE uuid = $5
		RETURNING uuid
	`

	rows, queryErr := r.db.Exec(query, table.Name, columnsJSON, time.Now(), table.UpdatedBy, table.Uuid)
	if queryErr != nil {
		return nil, utils.FormatError(queryErr, "update", utils.GetMethodName())
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return nil, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	if rowsAffected != 1 {
		return nil, errs.NewNotFoundError("table.error.notFound")
	}

	return table, nil
}

func (r *CoreTableRepository) Rename(tableUUID uuid.UUID, name string, authUserID uuid.UUID) (models.Table, error) {
	query := `
		UPDATE fluxton.tables 
		SET name = $1, updated_at = $2, updated_by = $3
		WHERE uuid = $4`

	rows, err := r.db.Exec(query, name, time.Now(), authUserID, tableUUID)
	if err != nil {
		return models.Table{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return models.Table{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	if rowsAffected != 1 {
		return models.Table{}, errs.NewNotFoundError("table.error.notFound")
	}

	return r.GetByID(tableUUID)
}

func (r *CoreTableRepository) Delete(tableUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.tables WHERE uuid = $1"
	res, err := r.db.Exec(query, tableUUID)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
