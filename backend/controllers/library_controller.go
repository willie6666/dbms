package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

// GetLibrary handles GET /api/protected/library
func GetLibrary(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var licenses []models.GameLicense
	// Preload Game and Game.Media details, and only fetch ACTIVE licenses
	if err := database.DB.Preload("Game").Preload("Game.Media").Where("user_id = ? AND status = ?", userID, "ACTIVE").Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch library"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": licenses})
}

// GetWishlist handles GET /api/protected/wishlist
func GetWishlist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var wishlist []models.WishList
	if err := database.DB.Preload("Game").Preload("Game.Media").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wishlist})
}

// AddToWishlist handles POST /api/protected/wishlist
func AddToWishlist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var input struct {
		GameID uint `json:"game_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.WishList
	if err := database.DB.Where("user_id = ? AND game_id = ?", userID, input.GameID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already in wishlist", "already_exists": true})
		return
	}

	var game models.Game
	if err := database.DB.First(&game, input.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	// Removed TAKEN_DOWN check per user requirement: Wishlist is a simple mark state.

	wishItem := models.WishList{
		UserID: userID,
		GameID: input.GameID,
	}

	if err := database.DB.Create(&wishItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist (might already exist)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to wishlist"})
}

// RemoveFromWishlist handles DELETE /api/protected/wishlist/:game_id
func RemoveFromWishlist(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("game_id")

	if err := database.DB.Where("user_id = ? AND game_id = ?", userID, gameID).Delete(&models.WishList{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from wishlist"})
}

// PlayGame handles GET /api/protected/library/:game_id/play
func PlayGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("game_id")

	var license models.GameLicense
	if err := database.DB.Where("user_id = ? AND game_id = ? AND status = ?", userID, gameID, "ACTIVE").First(&license).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this game or the license is inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game launched successfully", "auth_token": "mock-play-token-12345"})
}

// DownloadGame handles GET /api/protected/library/:game_id/download
func DownloadGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("game_id")

	var license models.GameLicense
	if err := database.DB.Where("user_id = ? AND game_id = ? AND status = ?", userID, gameID, "ACTIVE").First(&license).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this game or the license is inactive"})
		return
	}

	var gameFile models.GameMedia
	if err := database.DB.Where("game_id = ? AND media_type = ?", gameID, "game_file").First(&gameFile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No downloadable game file is available"})
		return
	}

	// Support both new path format (/downloads/{game_id}/{file}) and legacy (/downloads/{file})
	var fullPath string
	fileURL := gameFile.FileURL
	if strings.HasPrefix(fileURL, "/downloads/") {
		rel := strings.TrimPrefix(fileURL, "/downloads/")
		rel = filepath.Clean(rel)
		if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game file path"})
			return
		}
		fullPath = filepath.Join("assets", "game-files", rel)

		// Fallback for legacy DB entries: If file isn't found at the root, check inside the game_id subfolder
		if _, err := os.Stat(fullPath); os.IsNotExist(err) && !strings.Contains(rel, "/") && !strings.Contains(rel, "\\") {
			fallbackPath := filepath.Join("assets", "game-files", gameID, rel)
			if _, err2 := os.Stat(fallbackPath); err2 == nil {
				fullPath = fallbackPath
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game file path"})
		return
	}

	c.FileAttachment(fullPath, filepath.Base(fullPath))
}
