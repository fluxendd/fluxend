package repositories

import (
	"database/sql"
	"errors"
	"fluxton/models"
	"fluxton/pkg"
	flxErrs "fluxton/pkg/errors"
	"fluxton/requests"
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

func (r *ProjectRepository) ListForUser(paginationParams requests.PaginationParams, authUserId uuid.UUID) ([]models.Project, error) {
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

	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[models.Project]("projects"))

	params := map[string]interface{}{
		"user_uuid": authUserId,
		"sort":      paginationParams.Sort,
		"limit":     paginationParams.Limit,
		"offset":    offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var organization models.Project
		if err := rows.StructScan(&organization); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		projects = append(projects, organization)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return projects, nil
}

func (r *ProjectRepository) List(paginationParams requests.PaginationParams) ([]models.Project, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit
	query := `SELECT %s FROM fluxton.projects ORDER BY :sort DESC LIMIT :limit OFFSET :offset;`

	query = fmt.Sprintf(query, pkg.GetColumns[models.Project]())

	params := map[string]interface{}{
		"sort":   paginationParams.Sort,
		"limit":  paginationParams.Limit,
		"offset": offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.StructScan(&project); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return projects, nil
}

func (r *ProjectRepository) GetByUUID(projectUUID uuid.UUID) (models.Project, error) {
	query := "SELECT %s FROM fluxton.projects WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[models.Project]())

	var project models.Project
	err := r.db.Get(&project, query, projectUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, flxErrs.NewNotFoundError("project.error.notFound")
		}

		return models.Project{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) GetDatabaseNameByUUID(projectUUID uuid.UUID) (string, error) {
	query := "SELECT db_name FROM fluxton.projects WHERE uuid = $1"

	var dbName string
	err := r.db.Get(&dbName, query, projectUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", flxErrs.NewNotFoundError("project.error.notFound")
		}

		return "", pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return dbName, nil
}

func (r *ProjectRepository) GetUUIDByDatabaseName(dbName string) (uuid.UUID, error) {
	query := "SELECT uuid FROM fluxton.projects WHERE db_name = $1"

	var projectUUID uuid.UUID
	err := r.db.Get(&projectUUID, query, dbName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, flxErrs.NewNotFoundError("project.error.notFound")
		}

		return uuid.UUID{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return projectUUID, nil
}

func (r *ProjectRepository) GetOrganizationUUIDByProjectUUID(id uuid.UUID) (uuid.UUID, error) {
	query := "SELECT organization_uuid FROM fluxton.projects WHERE uuid = $1"

	var organizationUUID uuid.UUID
	err := r.db.Get(&organizationUUID, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, flxErrs.NewNotFoundError("project.error.notFound")
		}

		return uuid.UUID{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return organizationUUID, nil
}

func (r *ProjectRepository) ExistsByUUID(id uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE uuid = $1)"

	var exists bool
	err := r.db.Get(&exists, query, id)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *ProjectRepository) ExistsByNameForOrganization(name string, organizationUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.projects WHERE name = $1 AND organization_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, name, organizationUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *ProjectRepository) Create(project *models.Project) (*models.Project, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	query := `
		INSERT INTO fluxton.projects (
			name, db_name, description, db_port, 
			organization_uuid, created_by, updated_by
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING uuid
	`

	queryErr := tx.QueryRowx(
		query,
		project.Name,
		project.DBName,
		project.Description,
		project.DBPort,
		project.OrganizationUuid,
		project.CreatedBy,
		project.UpdatedBy,
	).Scan(&project.Uuid)

	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, pkg.FormatError(queryErr, "insert", pkg.GetMethodName())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) Update(project *models.Project) (*models.Project, error) {
	query := `
		UPDATE fluxton.projects 
		SET name = :name, description = :description, updated_at = :updated_at, updated_by = :updated_by
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, project)
	if err != nil {
		return &models.Project{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &models.Project{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return project, nil
}

func (r *ProjectRepository) UpdateStatusByDatabaseName(databaseName, status string) (bool, error) {
	query := "UPDATE fluxton.projects SET status = $1 WHERE db_name = $2"
	res, err := r.db.Exec(query, status, databaseName)
	if err != nil {
		return false, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}

func (r *ProjectRepository) Delete(projectUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.projects WHERE uuid = $1"
	res, err := r.db.Exec(query, projectUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}
