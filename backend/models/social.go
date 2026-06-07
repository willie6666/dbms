package models

import "time"

type Review struct {
	ReviewID  uint      `gorm:"primaryKey;column:review_id" json:"review_id"`
	GameID    uint      `gorm:"not null" json:"game_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"not null" json:"content"`
	Attitude  string    `gorm:"not null" json:"attitude"` // POSITIVE or NEGATIVE
	Status    string    `gorm:"default:VISIBLE" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type ReviewReply struct {
	ReviewReplyID uint   `gorm:"primaryKey;column:review_reply_id" json:"review_reply_id"`
	ReviewID      uint   `gorm:"not null;column:review_id" json:"review_id"`
	UserID        uint   `gorm:"not null;column:user_id" json:"user_id"`
	ParentReplyID *uint  `gorm:"column:parent_reply_id" json:"parent_reply_id"` // Nullable for direct replies to review
	Content       string `gorm:"not null;column:content" json:"content"`
	Status        string    `gorm:"default:VISIBLE;column:status" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Friendship struct {
	FriendshipID uint      `gorm:"primaryKey;column:friendship_id" json:"friendship_id"`
	SenderID     uint      `gorm:"not null" json:"sender_id"`
	ReceiverID   uint      `gorm:"not null" json:"receiver_id"`
	Status       string    `gorm:"default:PENDING" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Blacklist struct {
	BlacklistID uint      `gorm:"primaryKey;column:blacklist_id" json:"blacklist_id"`
	BlockerID   uint      `gorm:"not null" json:"blocker_id"`
	BlockedID   uint      `gorm:"not null" json:"blocked_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Message struct {
	MessageID  uint      `gorm:"primaryKey;column:message_id" json:"message_id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"not null" json:"content"`
	SentAt     time.Time `gorm:"autoCreateTime;column:sent_at" json:"sent_at"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
}
