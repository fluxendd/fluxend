package policies

import (
	"github.com/google/uuid"
)

func CanUpdateUser(userID, authenticatedUserId uuid.UUID) bool {
	return userID == authenticatedUserId
}
