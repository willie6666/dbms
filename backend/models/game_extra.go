package models

type Tag struct {
	TagID   uint   `gorm:"primaryKey;column:tag_id" json:"tag_id"`
	TagName string `gorm:"unique;not null;column:tag_name" json:"tag_name"`
}

type GameTag struct {
	GameID uint `gorm:"primaryKey;column:game_id" json:"game_id"`
	TagID  uint `gorm:"primaryKey;column:tag_id" json:"tag_id"`
}

type GameMedia struct {
	MediaID   uint   `gorm:"primaryKey;column:media_id" json:"media_id"`
	GameID    uint   `gorm:"not null;column:game_id" json:"game_id"`
	FileURL   string `gorm:"not null;column:file_url" json:"file_url"`
	MediaType string `gorm:"default:media;column:media_type" json:"media_type"`
}
