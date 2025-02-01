package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
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

	query = fmt.Sprintf(query, modelSkeleton.GetFieldsWithAlias(modelSkeleton.GetTableName()))

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

func (r *CoreTableRepository) GetByID(id uint) (models.Table, error) {
	query := "SELECT %s FROM tables WHERE id = $1"
	query = fmt.Sprintf(query, models.Table{}.GetFields())

	var table models.Table
	row := r.db.QueryRow(query, id)

	err := row.Scan(&table.ID, &table.ProjectID, &table.Name, &table.Fields, &table.CreatedAt, &table.UpdatedAt)
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

func (r *CoreTableRepository) Create(table *models.Table) (*models.Table, error) {
	fieldsJSON, err := table.MarshalJSONFields()
	if err != nil {
		return nil, fmt.Errorf("could not marshal fields: %v", err)
	}

	query := "INSERT INTO tables (name, project_id, fields) VALUES ($1, $2, $3) RETURNING id"
	queryErr := r.db.QueryRowx(query, table.Name, table.ProjectID, fieldsJSON).Scan(&table.ID)
	if queryErr != nil {
		return nil, fmt.Errorf("could not create table: %v", queryErr)
	}

	return table, nil
}

func (r *CoreTableRepository) Rename(tableID uint, name string) (models.Table, error) {
	query := `
		UPDATE tables 
		SET name = $1, updated_at = $2 
		WHERE id = $3`

	_, queryErr := r.db.Exec(query, name, time.Now(), tableID)
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
