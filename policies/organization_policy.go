package policies

type OrganizationPolicy interface {
	CanCreate(authenticatedUserId uint) bool
	CanUpdate(userID, authenticatedUserId uint) bool
}

type OrganizationPolicyImpl struct {
}

func NewOrganizationPolicy() OrganizationPolicy {
	return &OrganizationPolicyImpl{}
}

func (s *OrganizationPolicyImpl) CanCreate(authenticatedUserId uint) bool {
	return true
}

func (s *OrganizationPolicyImpl) CanView(organizationUserId, authenticatedUserId uint) bool {
	return organizationUserId == authenticatedUserId
}

func (s *OrganizationPolicyImpl) CanUpdate(organizationUserId, authenticatedUserId uint) bool {
	return organizationUserId == authenticatedUserId
}
