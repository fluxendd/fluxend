package configs

var messages = map[string]string{
	"auth.error.tokenRequired":      "Token is required",
	"auth.error.tokenInvalid":       "Invalid token provided",
	"user.error.notFound":           "User not found",
	"user.error.invalidCredentials": "Invalid credentials provided",
	"user.error.updateForbidden":    "You don't have permission to update user",
	"user.error.unauthenticated":    "Unauthenticated",
}

func Message(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}

	return key
}
