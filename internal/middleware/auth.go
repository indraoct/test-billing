package middleware

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ContextKey string

const (
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
)

// AuthMiddleware validates the Authorization header and attaches user info to the request context
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Authorization header is missing",
				})
			}

			// Check if the header is in the correct format (e.g., "Bearer <token>")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Invalid Authorization header format",
				})
			}

			token := parts[1]

			// Validate the token (replace with actual validation logic)
			userID, userRole, err := validateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Invalid token: " + err.Error(),
				})
			}

			// Attach user info to the context
			ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			c.Request().Context().Value(ctx)

			// Proceed to the next handler
			return next(c)
		}
	}
}

// validateToken is a placeholder function for token validation
func validateToken(token string) (int, string, error) {
	// Replace this with actual token validation logic, such as decoding a JWT
	if token == "admin-token" {
		return 1, "admin", nil // Example: valid admin token (just for example)
	}
	if token == "user-token" {
		return 1, "user", nil // Example: valid user token (just for example)
	}
	return 0, "", errors.New("token is not valid")
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// GetUserRoleFromContext retrieves the user role from the context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	userRole, ok := ctx.Value(UserRoleKey).(string)
	return userRole, ok
}
