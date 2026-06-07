package models

import "time"

// GameLicense corresponds to the 'game_licenses' table (Owned Games)
type GameLicense struct {
	LicenseID         uint      `gorm:"primaryKey;column:license_id" json:"license_id"`
	GameID            uint      `gorm:"not null" json:"game_id"`
	UserID            uint      `gorm:"not null" json:"user_id"`
	TransactionItemID uint      `gorm:"not null" json:"transaction_item_id"`
	AcquiredDate      time.Time `gorm:"autoCreateTime" json:"acquired_date"`
	Status            string    `gorm:"default:ACTIVE" json:"status"`

	// Preload relationship to fetch the actual Game details
	Game Game `gorm:"foreignKey:GameID;references:GameID" json:"game"`
}

// WishList corresponds to the 'wish_lists' table
type WishList struct {
	WishlistID uint      `gorm:"primaryKey;column:wishlist_id" json:"wishlist_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	GameID     uint      `gorm:"not null" json:"game_id"`
	AddedAt    time.Time `gorm:"autoCreateTime" json:"added_at"`

	// Preload relationship to fetch the actual Game details
	Game Game `gorm:"foreignKey:GameID;references:GameID" json:"game"`
}
