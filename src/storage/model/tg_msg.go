package model

type TgMsg struct {
	InternalID       int64  `gorm:"primaryKey;autoIncrement" json:"internal_id"`
	TgID             int64  `gorm:"not null;default:0" json:"tg_id"`
	TgUserID         int64  `gorm:"not null;default:0" json:"tg_user_id"`
	ChatID           int64  `gorm:"not null;default:0" json:"chat_id"`
	ReplyToMessageID int64  `gorm:"not null;default:0" json:"reply_to_message_id"`
	IsOutgoing       byte   `gorm:"not null;default:0" json:"is_outgoing"`
	Text             string `gorm:"not null" json:"text"`
	Date             int64  `gorm:"default:0" json:"date"`
	AddedTimestamp   int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}
