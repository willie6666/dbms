package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"
	"vapor_auror_backend/utils"

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
	tag := c.Query("tag")
	developer := c.Query("developer")
	sort := c.DefaultQuery("sort", "price_asc")

	query := database.DB.Model(&models.Game{}).Group("games.game_id").Where("games.status = ?", "ACTIVE")
	if q != "" {
		keyword := "%" + q + "%"
		query = query.
			Joins("LEFT JOIN game_tags q_game_tags ON q_game_tags.game_id = games.game_id").
			Joins("LEFT JOIN tags q_tags ON q_tags.tag_id = q_game_tags.tag_id").
			Joins("LEFT JOIN users q_developers ON q_developers.user_id = games.developer_id").
			Where("games.title ILIKE ? OR games.description ILIKE ? OR q_tags.tag_name ILIKE ? OR q_developers.username ILIKE ?", keyword, keyword, keyword, keyword)
	}
	if tag != "" {
		query = query.
			Joins("JOIN game_tags filter_game_tags ON filter_game_tags.game_id = games.game_id").
			Joins("JOIN tags filter_tags ON filter_tags.tag_id = filter_game_tags.tag_id").
			Where("filter_tags.tag_name ILIKE ?", tag)
	}
	if developer != "" {
		query = query.
			Joins("JOIN users filter_developers ON filter_developers.user_id = games.developer_id").
			Where("filter_developers.username ILIKE ?", developer)
	}
	if minPrice, err := strconv.ParseFloat(c.Query("min_price"), 64); err == nil {
		query = query.Where("games.price >= ?", minPrice)
	}
	if maxPrice, err := strconv.ParseFloat(c.Query("max_price"), 64); err == nil {
		query = query.Where("games.price <= ?", maxPrice)
	}

	if c.Query("hide_owned") == "true" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				if claims, err := utils.ParseToken(parts[1]); err == nil {
					userID := uint(claims["user_id"].(float64))
					query = query.Where("NOT EXISTS (SELECT 1 FROM game_licenses WHERE game_licenses.game_id = games.game_id AND game_licenses.user_id = ? AND game_licenses.status = 'ACTIVE') AND games.developer_id != ?", userID, userID)
				}
			}
		}
	}

	switch sort {
	case "price_desc":
		query = query.Order("games.price DESC, games.game_id DESC")
	case "price_asc":
		fallthrough
	default:
		query = query.Order("games.price ASC, games.game_id DESC")
	}

	// Retrieve games from the database with their media
	if err := query.Preload("Media").Preload("Tags").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games"})
		return
	}
	for i := range games {
		attachRatingSummary(&games[i])
		var dev models.User
		if err := database.DB.Select("username").First(&dev, games[i].DeveloperID).Error; err == nil {
			games[i].DeveloperName = dev.Username
		} else {
			games[i].DeveloperName = "未知"
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": games})
}

// GetGameByID handles GET /api/games/:id (UC-03)
func GetGameByID(c *gin.Context) {
	var game models.Game
	gameID := c.Param("id")

	if err := database.DB.Preload("Tags").First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	attachRatingSummary(&game)

	var media []models.GameMedia
	database.DB.Where("game_id = ?", gameID).Find(&media)

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
			"tags":           game.Tags,
			"reviews":        fullReviews,
		},
	})
}
