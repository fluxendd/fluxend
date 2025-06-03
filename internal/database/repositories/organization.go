package repositories

import (
	"fluxend/internal/domain/organization"
	"fluxend/internal/domain/shared"
	"fluxend/internal/domain/user"
	"fluxend/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type OrganizationRepository struct {
	db shared.DB
}

func NewOrganizationRepository(injector *do.Injector) (organization.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &OrganizationRepository{db: db}, nil
}

func (r *OrganizationRepository) ListForUser(paginationParams shared.PaginationParams, authUserID uuid.UUID) ([]organization.Organization, error) {
	offset := (paginationParams.Page - 1) * paginationParams.Limit

	query := `
		SELECT 
			%s 
		FROM 
			fluxend.organizations organizations
		JOIN 
			fluxend.organization_members organization_members ON organizations.uuid = organization_members.organization_uuid
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

	var organizations []organization.Organization
	return organizations, r.db.SelectNamedList(&organizations, query, params)
}

func (r *OrganizationRepository) ListUsers(organizationUUID uuid.UUID) ([]user.User, error) {
	query := `
		SELECT 
			%s 
		FROM 
			authentication.users users
		JOIN 
			fluxend.organization_members organization_members ON users.uuid = organization_members.user_uuid
		WHERE 
			organization_members.organization_uuid = $1
	`

	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[user.User]("users"))

	var users []user.User
	return users, r.db.Select(&users, query, organizationUUID)
}

func (r *OrganizationRepository) GetUser(organizationUUID, userUUID uuid.UUID) (user.User, error) {
	query := `
		SELECT 
			%s 
		FROM 
			authentication.users users
		JOIN 
			fluxend.organization_members organization_members ON users.uuid = organization_members.user_uuid
		WHERE 
			organization_members.organization_uuid = $1 AND organization_members.user_uuid = $2
	`
	query = fmt.Sprintf(query, pkg.GetColumnsWithAlias[user.User]("users"))

	var currentUser user.User
	return currentUser, r.db.GetWithNotFound(&currentUser, "organization.error.userNotFound", query, organizationUUID, userUUID)
}

func (r *OrganizationRepository) CreateUser(organizationUUID, userUUID uuid.UUID) error {
	query := "INSERT INTO fluxend.organization_members (organization_uuid, user_uuid) VALUES ($1, $2)"
	_, err := r.db.Exec(query, organizationUUID, userUUID)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}

func (r *OrganizationRepository) DeleteUser(organizationUUID, userUUID uuid.UUID) error {
	return r.db.ExecWithErr("DELETE FROM fluxend.organization_members WHERE organization_uuid = $1 AND user_uuid = $2", organizationUUID, userUUID)
}

func (r *OrganizationRepository) GetByUUID(organizationUUID uuid.UUID) (organization.Organization, error) {
	query := "SELECT %s FROM fluxend.organizations WHERE uuid = $1"
	query = fmt.Sprintf(query, pkg.GetColumns[organization.Organization]())

	var fetchedOrganization organization.Organization
	return fetchedOrganization, r.db.GetWithNotFound(&fetchedOrganization, "organization.error.notFound", query, organizationUUID)
}

func (r *OrganizationRepository) ExistsByID(organizationUUID uuid.UUID) (bool, error) {
	return r.db.Exists("fluxend.organizations", "uuid = $1", organizationUUID)
}

func (r *OrganizationRepository) Create(organization *organization.Organization, authUserID uuid.UUID) (*organization.Organization, error) {
	return organization, r.db.WithTransaction(func(tx shared.Tx) error {
		// Insert into organizations table
		query := "INSERT INTO fluxend.organizations (name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING uuid"
		if err := tx.QueryRowx(query, organization.Name, organization.CreatedBy, organization.UpdatedBy).Scan(&organization.Uuid); err != nil {
			return fmt.Errorf("could not create organization: %v", err)
		}

		// Insert into organization_members pivot table
		if err := r.createOrganizationUser(tx, organization.Uuid, authUserID); err != nil {
			return fmt.Errorf("could not insert into pivot table: %v", err)
		}

		return nil
	})
}

func (r *OrganizationRepository) Update(organizationInput *organization.Organization) (*organization.Organization, error) {
	query := `
		UPDATE fluxend.organizations 
		SET name = :name, updated_at = :updated_at, updated_by = :updated_by 
		WHERE uuid = :uuid`

	err := r.db.ExecWithErr(query, organizationInput)

	return organizationInput, err
}

func (r *OrganizationRepository) Delete(organizationUUID uuid.UUID) (bool, error) {
	rowsAffected, err := r.db.ExecWithRowsAffected("DELETE FROM fluxend.organizations WHERE uuid = $1", organizationUUID)
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (r *OrganizationRepository) IsOrganizationMember(organizationUUID, authUserID uuid.UUID) (bool, error) {
	condition := "organization_uuid = $1 AND user_uuid = $2"
	return r.db.Exists("fluxend.organization_members", condition, organizationUUID, authUserID)
}

func (r *OrganizationRepository) createOrganizationUser(tx shared.Tx, organizationUUID, userId uuid.UUID) error {
	query := "INSERT INTO fluxend.organization_members (organization_uuid, user_uuid) VALUES ($1, $2)"
	_, err := tx.Exec(query, organizationUUID, userId)
	if err != nil {
		return fmt.Errorf("could not insert into pivot table: %v", err)
	}

	return nil
}
