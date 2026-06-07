package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

// GetGames handles GET /api/games (UC-01, UC-02)
func GetGames(c *gin.Context) {
	var games []models.Game
	q := c.Query("q")

	query := database.DB
	if q != "" {
		// PostgreSQL ILIKE is case-insensitive
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}

	// Retrieve games from the database with their media
	if err := query.Preload("Media").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": games})
}

// GetGameByID handles GET /api/games/:id (UC-03)
func GetGameByID(c *gin.Context) {
	var game models.Game
	gameID := c.Param("id")

	if err := database.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	var media []models.GameMedia
	database.DB.Where("game_id = ?", gameID).Find(&media)

	var gameTags []models.GameTag
	database.DB.Where("game_id = ?", gameID).Find(&gameTags)
	var tags []models.Tag
	for _, gt := range gameTags {
		var tag models.Tag
		if err := database.DB.First(&tag, gt.TagID).Error; err == nil {
			tags = append(tags, tag)
		}
	}

	var reviews []models.Review
	database.DB.Where("game_id = ?", gameID).Find(&reviews)

	type ReviewWithReplies struct {
		models.Review
		Replies []models.ReviewReply `json:"replies"`
	}
	var fullReviews []ReviewWithReplies

	for _, r := range reviews {
		var replies []models.ReviewReply
		database.DB.Where("review_id = ?", r.ReviewID).Find(&replies)
		fullReviews = append(fullReviews, ReviewWithReplies{Review: r, Replies: replies})
	}

	var dev models.User
	developerName := "未知"
	if err := database.DB.First(&dev, game.DeveloperID).Error; err == nil {
		developerName = dev.Username
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"game":           game,
			"developer_name": developerName,
			"media":          media,
			"tags":           tags,
			"reviews":        fullReviews,
		},
	})
}
