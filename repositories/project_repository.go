package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"myapp/errs"
	"myapp/models"
	"myapp/utils"
	"time"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(injector *do.Injector) (*ProjectRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &ProjectRepository{db: db}, nil
}

func (r *ProjectRepository) ListForUser(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Project, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	modelSkeleton := models.Project{}

	query := `
		SELECT 
			%s 
		FROM 
			projects
		JOIN 
			organization_users ON projects.id = organization_users.organization_id
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

func (r *ProjectRepository) GetByIDForUser(id, authenticatedUserId uint) (models.Project, error) {
	query := "SELECT %s FROM projects WHERE id = $1"
	query = fmt.Sprintf(query, models.Project{}.GetFields())

	var organization models.Project
	err := r.db.Get(&organization, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, errs.NewNotFoundError("organization.error.notFound")
		}

		return models.Project{}, fmt.Errorf("could not fetch row: %v", err)
	}

	return organization, nil
}

func (r *ProjectRepository) ExistsByID(id uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *ProjectRepository) ExistsByName(name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM projects WHERE name = $1)"

	var exists bool
	err := r.db.Get(&exists, query, name)
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

	// Insert into projects table
	query := "INSERT INTO projects (name, db_name, organization_id) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowx(query, project.Name, project.DBName, project.OrganizationID).Scan(&project.ID)
	if err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("could not create project: %v", err)
	}

	// Insert into project_users pivot table
	/*err = r.createOrganizationUser(tx, project.ID, authenticatedUserId)
	if err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("could not insert into pivot table: %v", err)
	}*/

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %v", err)
	}

	return project, nil
}

func (r *ProjectRepository) Update(id uint, organization *models.Project) (*models.Project, error) {
	organization.UpdatedAt = time.Now()
	organization.ID = id

	query := `
		UPDATE projects 
		SET name = :name, updated_at = :updated_at 
		WHERE id = :id`

	res, err := r.db.NamedExec(query, organization)
	if err != nil {
		return &models.Project{}, fmt.Errorf("could not update row: %v", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Project{}, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return organization, nil
}

func (r *ProjectRepository) Delete(organizationId uint) (bool, error) {
	query := "DELETE FROM projects WHERE id = $1"
	res, err := r.db.Exec(query, organizationId)
	if err != nil {
		return false, fmt.Errorf("could not delete row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("could not determine affected rows: %v", err)
	}

	return rowsAffected == 1, nil
}

func (r *ProjectRepository) IsOrganizationUser(organizationId, authenticatedUserId uint) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM organization_users WHERE organization_id = $1 AND user_id = $2)"

	var exists bool
	err := r.db.Get(&exists, query, organizationId, authenticatedUserId)
	if err != nil {
		return false, fmt.Errorf("could not fetch row: %v", err)
	}

	return exists, nil
}

func (r *ProjectRepository) createOrganizationUser(tx *sqlx.Tx, organizationId, userId uint) error {
	query := "INSERT INTO organization_users (organization_id, user_id) VALUES ($1, $2)"
	_, err := tx.Exec(query, organizationId, userId)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}
