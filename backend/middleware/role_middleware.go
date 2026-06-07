package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole is the second layer of defense. It checks if the decoded JWT role matches the requirement.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the role that was attached to the context by the AuthMiddleware
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		// Check if the user's role matches the required role
		// ADMIN is a superuser and can access anything
		if roleStr := role.(string); roleStr != requiredRole && roleStr != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		// If the role matches, allow them to proceed to the Controller
		c.Next()
	}
}
