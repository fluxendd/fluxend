package repositories

import (
	"database/sql"
	"errors"
	"fluxton/internal/domain/organization"
	"fluxton/internal/domain/shared"
	"fluxton/internal/domain/user"
	"fluxton/pkg"
	flxErrs "fluxton/pkg/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type OrganizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(injector *do.Injector) (organization.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &OrganizationRepository{db: db}, nil
}

func (r *OrganizationRepository) ListForUser(paginationParams shared.PaginationParams, authUserID uuid.UUID) ([]organization.Organization, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT 
			%s 
		FROM 
			fluxton.organizations organizations
		JOIN 
			fluxton.organization_members organization_members ON organizations.uuid = organization_members.organization_uuid
		WHERE 
			organization_members.user_uuid = :user_uuid
		ORDER BY 
			:sort DESC
		LIMIT 
			:limit 
		OFFSET 
			:offset;

	`

	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[organization.Organization]("organizations"))

	params := map[string]interface{}{
		"user_uuid": authUserID,
		"sort":      paginationParams.Sort,
		"limit":     paginationParams.Limit,
		"offset":    offset,
	}

	rows, err := r.db.NamedQuery(query, params)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var organizations []organization.Organization
	for rows.Next() {
		var currentOrganization organization.Organization
		if err := rows.StructScan(&currentOrganization); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		organizations = append(organizations, currentOrganization)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return organizations, nil
}

func (r *OrganizationRepository) ListUsers(organizationUUID uuid.UUID) ([]user.User, error) {
	query := `
		SELECT 
			%s 
		FROM 
			authentication.users users
		JOIN 
			fluxton.organization_members organization_members ON users.uuid = organization_members.user_uuid
		WHERE 
			organization_members.organization_uuid = $1
	`

	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[user.User]("users"))
	rows, err := r.db.Queryx(query, organizationUUID)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var currentUser user.User
		if err := rows.StructScan(&currentUser); err != nil {
			return nil, pkg.FormatError(err, "scan", pkg.GetMethodName())
		}
		users = append(users, currentUser)
	}

	if err := rows.Err(); err != nil {
		return nil, pkg.FormatError(err, "iterate", pkg.GetMethodName())
	}

	return users, nil
}

func (r *OrganizationRepository) GetUser(organizationUUID, userUUID uuid.UUID) (user.User, error) {
	query := `
		SELECT 
			%s 
		FROM 
			authentication.users users
		JOIN 
			fluxton.organization_members organization_members ON users.uuid = organization_members.user_uuid
		WHERE 
			organization_members.organization_uuid = $1 AND organization_members.user_uuid = $2
	`
	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[user.User]("users"))

	var currentUser user.User
	err := r.db.Get(&currentUser, query, organizationUUID, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, flxErrs.NewNotFoundError("organization.error.userNotFound")
		}

		return user.User{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return currentUser, nil
}

func (r *OrganizationRepository) CreateUser(organizationUUID, userUUID uuid.UUID) error {
	query := "INSERT INTO fluxton.organization_members (organization_uuid, user_uuid) VALUES ($1, $2)"
	_, err := r.db.Exec(query, organizationUUID, userUUID)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}

func (r *OrganizationRepository) DeleteUser(organizationUUID, userUUID uuid.UUID) error {
	query := "DELETE FROM fluxton.organization_members WHERE organization_uuid = $1 AND user_uuid = $2"
	_, err := r.db.Exec(query, organizationUUID, userUUID)
	if err != nil {
		return pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	return nil
}

func (r *OrganizationRepository) GetByUUID(organizationUUID uuid.UUID) (organization.Organization, error) {
	query := "SELECT %s FROM fluxton.organizations WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[organization.Organization]())

	var fetchedOrganization organization.Organization
	err := r.db.Get(&fetchedOrganization, query, organizationUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return organization.Organization{}, flxErrs.NewNotFoundError("organization.error.notFound")
		}

		return organization.Organization{}, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return fetchedOrganization, nil
}

func (r *OrganizationRepository) ExistsByID(organizationUUID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.organizations WHERE uuid = $1)"
	var exists bool
	err := r.db.Get(&exists, query, organizationUUID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *OrganizationRepository) Create(organization *organization.Organization, authUserID uuid.UUID) (*organization.Organization, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}

	// Insert into organizations table
	query := "INSERT INTO fluxton.organizations (name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING uuid"
	queryErr := tx.QueryRowx(query, organization.Name, organization.CreatedBy, organization.UpdatedBy).Scan(&organization.Uuid)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not create organization: %v", queryErr)
	}

	// Insert into organization_members pivot table
	queryErr = r.createOrganizationUser(tx, organization.Uuid, authUserID)
	if queryErr != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("could not insert into pivot table: %v", queryErr)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return organization, nil
}

func (r *OrganizationRepository) Update(organizationInput *organization.Organization) (*organization.Organization, error) {
	query := `
		UPDATE fluxton.organizations 
		SET name = :name, updated_at = :updated_at, updated_by = :updated_by 
		WHERE uuid = :uuid`

	res, err := r.db.NamedExec(query, organizationInput)
	if err != nil {
		return &organization.Organization{}, pkg.FormatError(err, "update", pkg.GetMethodName())
	}

	_, err = res.RowsAffected()
	if err != nil {
		return &organization.Organization{}, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return organizationInput, nil
}

func (r *OrganizationRepository) Delete(organizationUUID uuid.UUID) (bool, error) {
	query := "DELETE FROM fluxton.organizations WHERE uuid = $1"
	res, err := r.db.Exec(query, organizationUUID)
	if err != nil {
		return false, pkg.FormatError(err, "delete", pkg.GetMethodName())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, pkg.FormatError(err, "affectedRows", pkg.GetMethodName())
	}

	return rowsAffected == 1, nil
}

func (r *OrganizationRepository) IsOrganizationMember(organizationUUID, authUserID uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM fluxton.organization_members WHERE organization_uuid = $1 AND user_uuid = $2)"

	var exists bool
	err := r.db.Get(&exists, query, organizationUUID, authUserID)
	if err != nil {
		return false, pkg.FormatError(err, "fetch", pkg.GetMethodName())
	}

	return exists, nil
}

func (r *OrganizationRepository) createOrganizationUser(tx *sqlx.Tx, organizationUUID, userId uuid.UUID) error {
	query := "INSERT INTO fluxton.organization_members (organization_uuid, user_uuid) VALUES ($1, $2)"
	_, err := tx.Exec(query, organizationUUID, userId)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}
