package controllers

import (
	"net/http"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
)

type AddCartInput struct {
	GameID uint `json:"game_id" binding:"required"`
}

// GetCart handles GET /api/protected/cart
func GetCart(c *gin.Context) {
	// 1. Get user_id from the AuthMiddleware
	userID, _ := c.Get("user_id")

	// 2. Fetch all cart items for this user, and Preload the "Game" details
	var cartItems []models.ShoppingCart
	if err := database.DB.Preload("Game").Preload("Game.Media").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cartItems})
}

// AddToCart handles POST /api/protected/cart
func AddToCart(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64)) // JWT numbers are parsed as float64

	var input AddCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Check if game already in cart
	var existing models.ShoppingCart
	if err := database.DB.Where("user_id = ? AND game_id = ?", userID, input.GameID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game already in cart"})
		return
	}

	// 1.5 Check if the user already OWNS the game (has an ACTIVE license)
	var license models.GameLicense
	if err := database.DB.Where("user_id = ? AND game_id = ? AND status = ?", userID, input.GameID, "ACTIVE").First(&license).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already own this game"})
		return
	}

	// 1.8 Check if the game is TAKEN_DOWN
	var game models.Game
	if err := database.DB.First(&game, input.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	if game.Status != "ACTIVE" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This game is not available for purchase"})
		return
	}

	// 2. Add to cart
	cartItem := models.ShoppingCart{
		UserID: userID,
		GameID: input.GameID,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game added to cart successfully"})
}

// RemoveFromCart handles DELETE /api/protected/cart/:game_id
func RemoveFromCart(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))
	gameID := c.Param("game_id")

	if err := database.DB.Where("user_id = ? AND game_id = ?", userID, gameID).Delete(&models.ShoppingCart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game removed from cart"})
}
