package models

import "time"

// Transaction corresponds to the 'transactions' table
type Transaction struct {
	TransactionID   uint              `gorm:"primaryKey;column:transaction_id" json:"transaction_id"`
	UserID          uint              `gorm:"not null" json:"user_id"`
	TotalAmount     float64           `gorm:"not null;default:0.00" json:"total_amount"`
	TransactionDate time.Time         `gorm:"autoCreateTime" json:"transaction_date"`
	ReceiptNumber   string            `gorm:"unique;not null" json:"receipt_number"`
	Items           []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}

// TransactionItem corresponds to the 'transaction_items' table
type TransactionItem struct {
	ItemID        uint    `gorm:"primaryKey;column:item_id" json:"item_id"`
	TransactionID uint    `gorm:"not null" json:"transaction_id"`
	GameID        uint    `gorm:"not null" json:"game_id"`
	PurchasePrice float64 `gorm:"not null;default:0.00" json:"purchase_price"`
	RefundStatus  string  `gorm:"-" json:"refund_status,omitempty"`

	// Preload relationship to fetch the actual Game details
	Game Game `gorm:"foreignKey:GameID;references:GameID" json:"game"`
}
