package models

import "time"

// ShoppingCart corresponds to the 'shopping_carts' table
type ShoppingCart struct {
	CartID  uint      `gorm:"primaryKey;column:cart_id" json:"cart_id"`
	UserID  uint      `gorm:"not null" json:"user_id"`
	GameID  uint      `gorm:"not null" json:"game_id"`
	AddedAt time.Time `gorm:"autoCreateTime" json:"added_at"`

	// Preload relationship to fetch the actual Game
	Game Game `gorm:"foreignKey:GameID;references:GameID" json:"game"`
}
