package middleware

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/responses"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return responses.UnauthorizedResponse(c, "auth.error.tokenRequired")
		}

		// Token usually comes in the format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // If the token doesn't start with "Bearer "
			return responses.UnauthorizedResponse(c, "auth.error.tokenInvalid")
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
			return responses.ErrorResponse(c, errs.NewUnauthorizedError("auth.error.tokenInvalid"))
		}

		// Optionally: You can store the user data from the token for later use in the request context
		// For example, store user ID or role in context for further authorization checks.
		userUUID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			return responses.ErrorResponse(c, errs.NewUnauthorizedError("auth.error.tokenInvalid"))
		}

		c.Set("user", models.AuthenticatedUser{
			ID:     userUUID,
			RoleID: int(claims["role_id"].(float64)),
		})

		// Proceed to the next handler if everything is valid
		return next(c)
	}
}
