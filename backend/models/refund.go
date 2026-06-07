package models

import "time"

// RefundRequest corresponds to the 'refund_requests' table
type RefundRequest struct {
	RefundID          uint      `gorm:"primaryKey;column:refund_id" json:"refund_id"`
	BuyerID           uint      `gorm:"not null" json:"buyer_id"`
	TransactionItemID uint      `gorm:"not null" json:"transaction_item_id"`
	HandledBy         *uint     `json:"handled_by"` // Pointer because it can be null
	Reason            string    `gorm:"not null" json:"reason"`
	RejectReason      string    `json:"reject_reason"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	ResolvedAt        *time.Time `json:"resolved_at"`
	Status            string    `gorm:"default:PENDING" json:"status"` // PENDING, APPROVED, REJECTED
}
