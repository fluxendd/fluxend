package repositories

import (
	"database/sql"
	"errors"
	"fluxton/errs"
	"fluxton/models"
	"fluxton/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"time"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(injector *do.Injector) (*ProjectRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &ProjectRepository{db: db}, nil
}

func (r *ProjectRepository) ListForUser(paginationParams utils.PaginationParams, authenticatedUserId uuid.UUID) ([]models.Project, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	modelSkeleton := models.Project{}

	query := `
		SELECT 
			%s 
		FROM 
			fluxton.projects projects
		JOIN 
			fluxton.organization_users organization_users ON projects.organization_id = organization_users.organization_id
		WHERE 
			organization_users.user_id = :user_id
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, modelSkeleton.GetColumnsWithAlias(modelSkeleton.GetTableName()))

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

func (r *ProjectRepository) GetByID(id uuid.UUID) (models.Project, error) {
	query := "SELECT %s FROM fluxton.projects WHERE id = $1"
	query = fmt.Sprintf(query, models.Project{}.GetColumns())

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

func (r *ProjectRepository) ExistsByID(id uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE id = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *ProjectRepository) ExistsByNameForOrganization(name string, organizationId uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE name = $1 AND organization_id = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, organizationId)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *ProjectRepository) Create(project *models.Project) (*models.Project, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	query := "INSERT INTO fluxton.projects (name, db_name, organization_id) VALUES ($1, $2, $3) RETURNING id"
	queryErr := tx.QueryRowx(query, project.Name, project.DBName, project.OrganizationID).Scan(&project.ID)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create project: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return project, nil
}

func (r *ProjectRepository) Update(id uuid.UUID, project *models.Project) (*models.Project, error) {
	project.UpdatedAt = time.Now()
	project.ID = id

	query := `
		UPDATE fluxton.projects 
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

func (r *ProjectRepository) Delete(projectId uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.projects WHERE id = $1"
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
