package model

import "time"

type WebSession struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	WebUserID  int       `gorm:"not null" json:"web_user_id"`
	SessionKey string    `gorm:"not null;unique" json:"session_key"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
}
