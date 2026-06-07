package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

type UploadGameInput struct {
	Title string  `json:"title" binding:"required"`
	Price float64 `json:"price" binding:"min=0"`
}

// GetDeveloperGames handles GET /api/developer/games
func GetDeveloperGames(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))

	var games []models.Game
	query := database.DB.Preload("Media")
	if role, _ := c.Get("role"); role != "ADMIN" {
		query = query.Where("developer_id = ?", developerID)
	}

	if err := query.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch developer games"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": games})
}

// UploadGame handles POST /api/developer/games
func UploadGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))

	var input UploadGameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game := models.Game{
		DeveloperID: developerID,
		Title:       input.Title,
		Price:       input.Price,
	}

	// Insert new game into the database
	if err := database.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload game"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Game uploaded successfully",
		"game":    game,
	})
}

// DeleteGame handles DELETE /api/developer/games/:id
func DeleteGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	// 1. Find the game first
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// 2. IMPORTANT SECURITY CHECK: Ensure the developer trying to delete it is the owner
	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only delete your own games"})
		return
	}

	// 3. Delete the game
	if err := database.DB.Delete(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

// UploadMedia handles POST /api/protected/developer/games/:id/media
func UploadMedia(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only upload media for your own games"})
		return
	}

	var input struct {
		FileURL   string `json:"file_url" binding:"required"`
		MediaType string `json:"media_type" binding:"required"` // 'media' or 'game_file'
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media := models.GameMedia{
		GameID:    game.GameID,
		FileURL:   input.FileURL,
		MediaType: input.MediaType,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload media"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Media uploaded successfully", "data": media})
}

// DeleteMedia handles DELETE /api/protected/developer/games/:id/media/:media_id
func DeleteMedia(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")
	mediaID := c.Param("media_id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only manage your own games"})
		return
	}

	if err := database.DB.Where("media_id = ? AND game_id = ?", mediaID, game.GameID).Delete(&models.GameMedia{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

// GetGameStats handles GET /api/protected/developer/games/:id/stats
func GetGameStats(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only view stats for your own games"})
		return
	}

	// Advanced Aggregation Query (SUM and COUNT)
	var result struct {
		TotalSales   int     `json:"total_sales"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	database.DB.Table("transaction_items").
		Select("count(*) as total_sales, COALESCE(sum(purchase_price), 0) as total_revenue").
		Where("game_id = ?", game.GameID).
		Scan(&result)

	c.JSON(http.StatusOK, gin.H{"stats": result})
}

// GetTags handles GET /api/tags (Public)
func GetTags(c *gin.Context) {
	var tags []models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag handles POST /api/tags (Developer)
func CreateTag(c *gin.Context) {
	var input struct {
		TagName string `json:"tag_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := models.Tag{TagName: input.TagName}
	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag (might already exist)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tag created successfully", "data": tag})
}

// AddTagToGame handles POST /api/protected/developer/games/:id/tags
func AddTagToGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Not your game"})
		return
	}

	var input struct {
		TagID uint `json:"tag_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gameTag := models.GameTag{GameID: game.GameID, TagID: input.TagID}
	if err := database.DB.Create(&gameTag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added to game"})
}

// RemoveTagFromGame handles DELETE /api/protected/developer/games/:id/tags/:tag_id
func RemoveTagFromGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")
	tagID := c.Param("tag_id")

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := database.DB.Where("game_id = ? AND tag_id = ?", game.GameID, tagID).Delete(&models.GameTag{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed from game"})
}
