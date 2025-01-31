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

func (r *CoreTableRepository) ListForUser(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Project, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	modelSkeleton := models.Project{}

	query := `
		SELECT 
			%s 
		FROM 
			projects
		JOIN 
			organization_users ON projects.organization_id = organization_users.organization_id
		WHERE 
			organization_users.user_id = :user_id
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, modelSkeleton.GetFieldsWithAlias(modelSkeleton.GetTableName()))

	params := map[string]interface{}{
		"user_id": authenticatedUserId,
		"sort":    paginationParams.Sort,
		"limit":   paginationParams.Limit,
		"offset":  offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var organization models.Project
		if err := rows.StructScan(&organization); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		projects = append(projects, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over rows: %v", err)
	}

	return projects, nil
}

func (r *CoreTableRepository) GetByID(id uint) (models.Project, error) {
	query := "SELECT %s FROM projects WHERE id = $1"
	query = fmt.Sprintf(query, models.Project{}.GetFields())

	var project models.Project
	err := r.db.Get(&project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, errs.NewNotFoundError("project.error.notFound")
		}

		return models.Project{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return project, nil
}

func (r *CoreTableRepository) ExistsByID(id uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *CoreTableRepository) ExistsByNameForOrganization(name string, organizationId uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM projects WHERE name = $1 AND organization_id = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, organizationId)
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

func (r *CoreTableRepository) Update(id uint, project *models.Project) (*models.Project, error) {
	project.UpdatedAt = time.Now()
	project.ID = id

	query := `
		UPDATE projects 
		SET name = :name, updated_at = :updated_at 
		WHERE id = :id`

	res, err := r.db.NamedExec(query, project)
	if err != nil {
		return &models.Project{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Project{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return project, nil
}

func (r *CoreTableRepository) Delete(projectId uint) (bool, error) {
	query := "DELETE FROM projects WHERE id = $1"
	res, err := r.db.Exec(query, projectId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}
