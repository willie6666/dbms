package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"
	"vapor_auror_backend/utils"

	"github.com/gin-gonic/gin"
)

// GetUsers handles GET /api/users
func GetUsers(c *gin.Context) {
	var users []models.User
	// Retrieve all users from the database
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// UpdateProfile handles PUT /api/users/profile
func UpdateProfile(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != "" && len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Username != "" && input.Username != user.Username {
		var existing models.User
		if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
			return
		}
		user.Username = input.Username
	}
	if input.Email != "" && input.Email != user.Email {
		var existing models.User
		if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
			return
		}
		user.Email = input.Email
	}
	if input.Password != "" {
		hash, err := utils.HashPassword(input.Password)
		if err == nil {
			user.PasswordHash = hash
		}
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}
