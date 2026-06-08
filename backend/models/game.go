package models

// Game corresponds to the 'games' table in the database
type Game struct {
	GameID        uint        `gorm:"primaryKey;column:game_id" json:"game_id"`
	DeveloperID   uint        `gorm:"not null" json:"developer_id"`
	Title         string      `gorm:"not null" json:"title"`
	Description   string      `gorm:"column:description" json:"desc"`
	Price         float64     `gorm:"not null;default:0.00" json:"price"`
	OverallRating float64     `gorm:"default:0.00" json:"overall_rating"`
	RatingCount   int64       `gorm:"-" json:"rating_count"`
	Media         []GameMedia `gorm:"foreignKey:GameID" json:"media,omitempty"`
	Tags          []Tag       `gorm:"many2many:game_tags;foreignKey:GameID;joinForeignKey:GameID;References:TagID;joinReferences:TagID" json:"tags,omitempty"`
	DeveloperName string      `gorm:"-" json:"developer_name"`
	Status        string      `gorm:"default:DRAFT" json:"status"`
}
