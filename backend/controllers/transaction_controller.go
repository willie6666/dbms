package controllers

import (
	"fmt"
	"net/http"
	"time"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Checkout handles POST /api/protected/checkout
func Checkout(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	// 1. Start a database transaction
	// This ensures that if anything fails (e.g., deducting money, giving license),
	// everything is rolled back, preventing half-finished transactions.
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch user's cart items with Game prices
		var cartItems []models.ShoppingCart
		if err := tx.Preload("Game").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return fmt.Errorf("Cart is empty")
		}

		// 2. Calculate Total Amount & Check Ownership
		var totalAmount float64 = 0
		for _, item := range cartItems {
			if item.Game.Status != "ACTIVE" {
				return fmt.Errorf("Game '%s' is no longer available for purchase", item.Game.Title)
			}
			// Double check ownership in case of concurrent requests or bypassed frontend
			var existingLicense models.GameLicense
			if err := tx.Where("user_id = ? AND game_id = ? AND status = ?", userID, item.GameID, "ACTIVE").First(&existingLicense).Error; err == nil {
				return fmt.Errorf("You already own '%s'", item.Game.Title)
			}
			totalAmount += item.Game.Price
		}

		// 3. Create Transaction Record
		receiptNumber := fmt.Sprintf("REC-%d-%d", userID, time.Now().UnixNano())
		transaction := models.Transaction{
			UserID:        userID,
			TotalAmount:   totalAmount,
			ReceiptNumber: receiptNumber,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// 4. Create Transaction Items & Game Licenses
		for _, cartItem := range cartItems {
			// 4a. Create TransactionItem
			txItem := models.TransactionItem{
				TransactionID: transaction.TransactionID,
				GameID:        cartItem.GameID,
				PurchasePrice: cartItem.Game.Price,
			}
			if err := tx.Create(&txItem).Error; err != nil {
				return err
			}

			// 4b. Create GameLicense (Grant access to the player)
			license := models.GameLicense{
				GameID:            cartItem.GameID,
				UserID:            userID,
				TransactionItemID: txItem.ItemID,
			}
			if err := tx.Create(&license).Error; err != nil {
				return err
			}
		}

		// 5. Clear the Shopping Cart
		if err := tx.Where("user_id = ?", userID).Delete(&models.ShoppingCart{}).Error; err != nil {
			return err
		}

		// Return nil commits the database transaction!
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Checkout failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful. Games added to your library!"})
}

// GetTransactions handles GET /api/protected/transactions
func GetTransactions(c *gin.Context) {
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	var transactions []models.Transaction
	// Preload Items, Game inside the Items, and Media inside the Game
	if err := database.DB.Preload("Items").Preload("Items.Game").Preload("Items.Game.Media").Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	for ti := range transactions {
		for ii := range transactions[ti].Items {
			var refund models.RefundRequest
			itemID := transactions[ti].Items[ii].ItemID
			if err := database.DB.Where("transaction_item_id = ?", itemID).Order("created_at desc").First(&refund).Error; err == nil {
				transactions[ti].Items[ii].RefundStatus = refund.Status
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
