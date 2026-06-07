package models

import "time"

// User corresponds to the 'users' table in the database
type User struct {
	UserID           uint      `gorm:"primaryKey;column:user_id" json:"id"`
	Username         string    `gorm:"unique;not null" json:"username"`
	Email            string    `gorm:"unique;not null" json:"email"`
	PasswordHash     string    `gorm:"not null" json:"-"`
	RegistrationDate time.Time `gorm:"autoCreateTime" json:"created_at"`
	LastVisitIP      *string   `json:"last_visit_ip"` // nullable, represents the VARCHAR(45) IP
	Role             string    `gorm:"not null" json:"role"`
	Status           string    `gorm:"default:OFFLINE" json:"status"`
	Permission       string    `gorm:"default:ACTIVE" json:"permission"`
}
