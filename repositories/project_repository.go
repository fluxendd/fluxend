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
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(injector *do.Injector) (*ProjectRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &ProjectRepository{db: db}, nil
}

func (r *ProjectRepository) ListForUser(paginationParams utils.PaginationParams, authUserId uuid.UUID) ([]models.Project, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `
		SELECT 
			%s 
		FROM 
			fluxton.projects projects
		JOIN 
			fluxton.organization_members organization_members ON projects.organization_uuid = organization_members.organization_uuid
		WHERE 
			organization_members.user_uuid = :user_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, utils.GetColumnsWithAlias[models.Project]("projects"))

	params := map[string]interface{}{
		"user_uuid": authUserId,
		"sort":      paginationParams.Sort,
		"limit":     paginationParams.Limit,
		"offset":    offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, utils.FormatError(err, "select", utils.GetMethodName())
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var organization models.Project
		if err := rows.StructScan(&organization); err != nil {
			return nil, utils.FormatError(err, "scan", utils.GetMethodName())
		}
		projects = append(projects, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.FormatError(err, "iterate", utils.GetMethodName())
	}

	return projects, nil
}

func (r *ProjectRepository) GetByUUID(projectUUID uuid.UUID) (models.Project, error) {
	query := "SELECT %s FROM fluxton.projects WHERE uuid = $1"
	query = fmt.Sprintf(query, utils.GetColumns[models.Project]())

	var project models.Project
	err := r.db.Get(&project, query, projectUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, errs.NewNotFoundError("project.error.notFound")
		}

		return models.Project{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) GetOrganizationUUIDByProjectUUID(id uuid.UUID) (uuid.UUID, error) {
	query := "SELECT organization_uuid FROM fluxton.projects WHERE uuid = $1"

	var organizationUUID uuid.UUID
	err := r.db.Get(&organizationUUID, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, errs.NewNotFoundError("project.error.notFound")
		}

		return uuid.UUID{}, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return organizationUUID, nil
}

func (r *ProjectRepository) ExistsByID(id uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *ProjectRepository) ExistsByNameForOrganization(name string, organizationUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE name = $1 AND organization_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, organizationUUID)
	if err != nil {
		return false, utils.FormatError(err, "fetch", utils.GetMethodName())
	}

	return exists, nil
}

func (r *ProjectRepository) Create(project *models.Project) (*models.Project, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, utils.FormatError(err, "transactionBegin", utils.GetMethodName())
	}

	query := "INSERT INTO fluxton.projects (name, db_name, db_port, organization_uuid, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING uuid"
	queryErr := tx.QueryRowx(query, project.Name, project.DBName, project.DBPort, project.OrganizationUuid, project.CreatedBy, project.UpdatedBy).Scan(&project.Uuid)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create project: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, utils.FormatError(err, "transactionCommit", utils.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) Update(project *models.Project) (*models.Project, error) {
	query := `
		UPDATE fluxton.projects 
		SET name = :name, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, project)
	if err != nil {
		return &models.Project{}, utils.FormatError(err, "update", utils.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Project{}, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) Delete(projectId uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.projects WHERE uuid = $1"
	res, err := r.db.Exec(query, projectId)
	if err != nil {
		return false, utils.FormatError(err, "delete", utils.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, utils.FormatError(err, "affectedRows", utils.GetMethodName())
	}

	return rowsAffected == 1, nil
}
