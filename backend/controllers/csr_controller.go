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
	var requests []models.RefundRequest
	if err := database.DB.Find(&requests).Error; err != nil {
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
