package controllers

import (
	"net/http"
	"strings"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"
	"vapor_auror_backend/utils"

	"github.com/gin-gonic/gin"
)

// Input structures for accepting JSON data from the frontend
type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	IsDeveloper bool   `json:"is_developer"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration (POST /api/auth/register)
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. Bind JSON from HTTP Request to our input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Minimal email format check: must contain "@"
	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// 2. Hash the password using Bcrypt
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Create the User model to save
	// Assign role based on IsDeveloper flag
	role := "USERS"
	if input.IsDeveloper {
		role = "DEVELOPER"
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         role,
		Status:       "ONLINE",
		Permission:   "ACTIVE",
	}

	// 4. Save to Database using GORM
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user (username or email might already exist)"})
		return
	}

	// 5. Generate JWT Token (Auto-login after registration)
	token, err := utils.GenerateToken(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"token":   token,
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login handles user authentication and JWT generation (POST /api/auth/login)
func Login(c *gin.Context) {
	var input LoginInput

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Minimal email format check: must contain "@"
	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// 2. Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. Compare hashed passwords
	if match := utils.CheckPassword(input.Password, user.PasswordHash); !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if user.Permission != "ACTIVE" || user.Role == "NULL" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This account is not active"})
		return
	}

	// 4. Generate JWT Token
	token, err := utils.GenerateToken(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Logout handles POST /api/auth/logout
func Logout(c *gin.Context) {
	// Since we use stateless JWTs, "logout" simply informs the client to drop the token.
	// If a token blacklist is needed, it would be implemented here.
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully. Please remove your token."})
}
