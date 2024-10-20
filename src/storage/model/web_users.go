package model

type WebUser struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username       string `gorm:"unique;not null" json:"username"`
	Password       string `gorm:"not null" json:"password"`
	Role           string `gorm:"not null" json:"role"`
	AddedTimestamp int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}
