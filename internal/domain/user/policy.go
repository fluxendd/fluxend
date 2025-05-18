package user

import (
	"github.com/google/uuid"
	"github.com/samber/do"
)

type Policy struct {
}

func NewUserPolicy(injector *do.Injector) (*Policy, error) {
	return &Policy{}, nil
}

func (s *Policy) CanUpdateUser(userID, authUserId uuid.UUID) bool {
	return userID == authUserId
}
