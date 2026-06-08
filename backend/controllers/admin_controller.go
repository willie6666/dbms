package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

// SuspendUser handles PUT /api/admin/users/:id/suspend
func SuspendUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Permission == "DELETED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify suspension status of a deleted user"})
		return
	}

	// Toggle suspension
	if user.Permission == "ACTIVE" {
		user.Permission = "DEACTIVE"
	} else {
		user.Permission = "ACTIVE"
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User account has been suspended"})
}

// DeleteUser handles DELETE /api/admin/users/:id
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Permission = "DELETED"
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Cascade soft-delete: if this user is a developer, take down all their games
	if user.Role == "DEVELOPER" {
		database.DB.Model(&models.Game{}).Where("developer_id = ?", user.UserID).Update("status", "TAKEN_DOWN")

		// Cascade revoke licenses for all games owned by this developer
		var developerGames []uint
		database.DB.Model(&models.Game{}).Where("developer_id = ?", user.UserID).Pluck("game_id", &developerGames)
		if len(developerGames) > 0 {
			database.DB.Model(&models.GameLicense{}).Where("game_id IN ?", developerGames).Update("status", "REVOKED")
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User completely removed"})
}

// ChangeUserRole handles PUT /api/admin/users/:id/role
func ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Update("role", input.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// AdminDeleteGame handles DELETE /api/admin/games/:id
func AdminDeleteGame(c *gin.Context) {
	gameID := c.Param("id")
	if err := database.DB.Model(&models.Game{}).Where("game_id = ?", gameID).Update("status", "TAKEN_DOWN").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to take down game"})
		return
	}

	// Cascade revoke all licenses for this game
	database.DB.Model(&models.GameLicense{}).Where("game_id = ?", gameID).Update("status", "REVOKED")

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully by Admin"})
}
