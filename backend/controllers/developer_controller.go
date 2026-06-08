package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

type UploadGameInput struct {
	Title string  `json:"title" binding:"required"`
	Price float64 `json:"price" binding:"min=0"`
	Desc  string  `json:"desc"`
}

type UpdateGameInput struct {
	Price float64 `json:"price" binding:"min=0"`
	Desc  string  `json:"desc"`
}

// GetDeveloperGames handles GET /api/developer/games
func GetDeveloperGames(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))

	var games []models.Game
	query := database.DB.Preload("Media").Preload("Tags")
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
		Description: input.Desc,
		Price:       input.Price,
		Status:      "DRAFT",
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

// PublishGame handles PUT /api/developer/games/:id/publish
func PublishGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var game models.Game
	// Preload Tags to check count
	if err := database.DB.Preload("Tags").First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if role, _ := c.Get("role"); role != "ADMIN" && game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only publish your own games"})
		return
	}

	if len(game.Tags) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game must have at least 1 tag to be published"})
		return
	}

	if err := database.DB.Model(&game).Update("status", "ACTIVE").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game published successfully"})
}

// UpdateGame handles PUT /api/developer/games/:id
func UpdateGame(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	developerID := uint(userIDFloat.(float64))
	gameID := c.Param("id")

	var input UpdateGameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	if role, _ := c.Get("role"); role != "ADMIN" && game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only edit your own games"})
		return
	}

	if err := database.DB.Model(&game).Updates(map[string]interface{}{
		"price":       input.Price,
		"description": input.Desc,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	database.DB.Preload("Media").First(&game, game.GameID)
	c.JSON(http.StatusOK, gin.H{"message": "Game updated successfully", "game": game})
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

	// 3. Delete the game (Soft delete by setting status to TAKEN_DOWN)
	if err := database.DB.Model(&game).Update("status", "TAKEN_DOWN").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to take down game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
}

// UploadMedia handles POST /api/developer/games/:id/media
// Accepts multipart/form-data with fields:
//   - file       : the file to upload
//   - media_type : "media" (image, default) or "game_file"
//
// Storage paths:
//
//	media     → assets/images/{game_id}/{sha256}.{ext}      served at /media/images/{game_id}/{sha256}.{ext}
//	game_file → assets/game-files/{game_id}/{original_name} served at /downloads/{game_id}/{original_name}
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

	// Parse the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file field"})
		return
	}

	mediaType := c.DefaultPostForm("media_type", "media")

	// Open and read file bytes
	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded file"})
		return
	}

	var dirPath, fileName, fileURL string

	if mediaType == "game_file" {
		// game_file: store under assets/game-files/{game_id}/{original_filename}
		// Sanitize the original filename (strip any directory components)
		originalName := filepath.Base(fileHeader.Filename)
		if originalName == "." || originalName == "" {
			originalName = "game_file.bin"
		}
		dirPath = filepath.Join("assets", "game-files", gameID)
		fileName = originalName
		fileURL = fmt.Sprintf("/downloads/%s/%s", gameID, fileName)
	} else {
		// media (image): store under assets/images/{game_id}/{sha256}.{ext}
		sum := sha256.Sum256(fileBytes)
		hashHex := fmt.Sprintf("%x", sum)
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext == "" {
			ext = ".bin"
		}
		dirPath = filepath.Join("assets", "images", gameID)
		fileName = hashHex + ext
		fileURL = fmt.Sprintf("/media/images/%s/%s", gameID, fileName)
	}

	// Ensure directory exists
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Write file to disk
	destPath := filepath.Join(dirPath, fileName)
	if err := os.WriteFile(destPath, fileBytes, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	media := models.GameMedia{
		GameID:    game.GameID,
		FileURL:   fileURL,
		MediaType: mediaType,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Media uploaded successfully",
		"data":     media,
		"file_url": fileURL,
	})
}

// DeleteMedia handles DELETE /api/developer/games/:id/media/:media_id
// Also removes the physical file from disk.
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

	// Fetch the media record first so we can remove the file
	var media models.GameMedia
	if err := database.DB.Where("media_id = ? AND game_id = ?", mediaID, game.GameID).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// Resolve the physical file path from the stored URL
	physicalPath := ""
	switch {
	case strings.HasPrefix(media.FileURL, "/media/images/"):
		// /media/images/{game_id}/{file} → assets/images/{game_id}/{file}
		rel := strings.TrimPrefix(media.FileURL, "/media/images/")
		physicalPath = filepath.Join("assets", "images", rel)
	case strings.HasPrefix(media.FileURL, "/downloads/"):
		// /downloads/{game_id}/{file} → assets/game-files/{game_id}/{file}
		rel := strings.TrimPrefix(media.FileURL, "/downloads/")
		physicalPath = filepath.Join("assets", "game-files", rel)
	}

	// Delete DB record
	if err := database.DB.Delete(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	// Best-effort physical file removal (ignore error — file may already be gone)
	if physicalPath != "" {
		_ = os.Remove(physicalPath)
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

	tagName := strings.TrimSpace(input.TagName)
	if tagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag name is required"})
		return
	}

	var existing models.Tag
	if err := database.DB.Where("LOWER(tag_name) = LOWER(?)", tagName).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Tag already exists", "data": existing})
		return
	}

	tag := models.Tag{TagName: tagName}
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

	if role, _ := c.Get("role"); role != "ADMIN" && game.DeveloperID != developerID {
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
	if err := database.DB.Where("game_id = ? AND tag_id = ?", game.GameID, input.TagID).FirstOrCreate(&gameTag).Error; err != nil {
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

	if role, _ := c.Get("role"); role != "ADMIN" && game.DeveloperID != developerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := database.DB.Where("game_id = ? AND tag_id = ?", game.GameID, tagID).Delete(&models.GameTag{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed from game"})
}
