package middleware

import (
	"net/http"
	"strings"
	"vapor_auror_backend/utils"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that checks if the request has a valid JWT token
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header (Format: "Bearer <token>")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// 2. Extract the token from the "Bearer " prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. Parse and validate the token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Attach the user_id and role to the Gin context so downstream controllers can use it
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		// 5. Proceed to the next handler (the actual API controller)
		c.Next()
	}
}
