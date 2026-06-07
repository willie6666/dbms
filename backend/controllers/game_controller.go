package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

func attachRatingSummary(game *models.Game) {
	var result struct {
		Rating float64
		Count  int64
	}
	database.DB.Table("reviews").
		Select("COALESCE(AVG(CASE WHEN attitude = 'POSITIVE' THEN 5.0 ELSE 1.0 END), 0) AS rating, COUNT(*) AS count").
		Where("game_id = ? AND status = ?", game.GameID, "VISIBLE").
		Scan(&result)
	game.OverallRating = result.Rating
	game.RatingCount = result.Count
}

func refreshStoredGameRating(gameID uint) {
	var game models.Game
	if err := database.DB.First(&game, gameID).Error; err != nil {
		return
	}
	attachRatingSummary(&game)
	database.DB.Model(&models.Game{}).Where("game_id = ?", gameID).Update("overall_rating", game.OverallRating)
}

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
	for i := range games {
		attachRatingSummary(&games[i])
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
	attachRatingSummary(&game)

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
