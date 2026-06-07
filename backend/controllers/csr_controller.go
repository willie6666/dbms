package controllers

import (
	"net/http"
	"time"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRefundRequests handles GET /api/csr/refunds
func GetRefundRequests(c *gin.Context) {
	type RefundDTO struct {
		RefundID          uint      `json:"refund_id"`
		BuyerID           uint      `json:"buyer_id"`
		Username          string    `json:"username"`
		TransactionItemID uint      `json:"transaction_item_id"`
		GameID            uint      `json:"game_id"`
		GameTitle         string    `json:"game_title"`
		Amount            float64   `json:"amount"`
		Reason            string    `json:"reason"`
		RejectReason      string    `json:"reject_reason"`
		Status            string    `json:"status"`
		CreatedAt         time.Time `json:"created_at"`
	}

	var requests []RefundDTO
	if err := database.DB.Table("refund_requests r").
		Select("r.refund_id, r.buyer_id, u.username, r.transaction_item_id, ti.game_id, g.title AS game_title, ti.purchase_price AS amount, r.reason, r.reject_reason, r.status, r.created_at").
		Joins("JOIN users u ON u.user_id = r.buyer_id").
		Joins("JOIN transaction_items ti ON ti.item_id = r.transaction_item_id").
		Joins("JOIN games g ON g.game_id = ti.game_id").
		Order("r.created_at DESC").
		Scan(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch refund requests"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// ProcessRefund handles PUT /api/csr/refunds/:id
func ProcessRefund(c *gin.Context) {
	csrIDFloat, _ := c.Get("user_id")
	csrID := uint(csrIDFloat.(float64))
	refundID := c.Param("id")

	var input struct {
		Status       string `json:"status" binding:"required"` // "APPROVED" or "REJECTED"
		RejectReason string `json:"reject_reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var request models.RefundRequest
		if err := tx.First(&request, refundID).Error; err != nil {
			return err
		}

		if request.Status != "PENDING" {
			return gorm.ErrInvalidData
		}

		now := time.Now()
		request.Status = input.Status
		request.HandledBy = &csrID
		request.ResolvedAt = &now
		if input.Status == "REJECTED" {
			request.RejectReason = input.RejectReason
		}

		if err := tx.Save(&request).Error; err != nil {
			return err
		}

		// If APPROVED, we must REVOKE the GameLicense
		if input.Status == "APPROVED" {
			// Find the license tied to this transaction item
			if err := tx.Model(&models.GameLicense{}).
				Where("transaction_item_id = ?", request.TransactionItemID).
				Update("status", "REVOKED").Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refund. Is it already processed?"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refund processed successfully"})
}
