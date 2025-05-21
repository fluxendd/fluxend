package middlewares

import (
	"fluxton/internal/api/response"
	"fluxton/internal/config/constants"
	"fluxton/internal/domain/auth"
	"fluxton/internal/domain/user"
	"fluxton/pkg/errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
)

func AuthMiddleware(userRepo user.Repository) echo.MiddlewareFunc {
	// Outer function accepts the next handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Inner function executes for each request
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.UnauthorizedResponse(c, "auth.error.tokenRequired")
			}

			// Token usually comes in the format "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader { // If the token doesn't start with "Bearer "
				return response.UnauthorizedResponse(c, "auth.error.bearerInvalid")
			}

			// Parse the token
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Check the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Return the JWT_SECRET as the key
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				// Token is invalid or expired
				return response.ErrorResponse(c, errors.NewUnauthorizedError("auth.error.tokenInvalid"))
			}

			userUUID, err := uuid.Parse(claims["uuid"].(string))
			if err != nil {
				return response.ErrorResponse(c, errors.NewUnauthorizedError("auth.error.tokenInvalid"))
			}

			loggedInJWTVersion := int(claims["version"].(float64))
			latestVersion, err := userRepo.GetJWTVersion(userUUID)
			if err != nil {
				return response.ErrorResponse(c, err)
			}

			// Allow a max 5 sessions to be active at the same time
			if (latestVersion - loggedInJWTVersion) >= constants.UserMaxLoginSessions {
				return response.UnauthorizedResponse(c, "auth.error.tokenInvalid")
			}

			c.Set("user", auth.User{
				Uuid:   userUUID,
				RoleID: int(claims["role_id"].(float64)),
			})

			// Proceed to the next handler if everything is valid
			return next(c)
		}
	}
}
