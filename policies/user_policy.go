package policies

import (
	"github.com/google/uuid"
)

func CanUpdateUser(userID, authUserId uuid.UUID) bool {
	return userID == authUserId
}
