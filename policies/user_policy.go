package policies

func CanUpdateUser(userID, authenticatedUserId uint) bool {
	return userID == authenticatedUserId
}
