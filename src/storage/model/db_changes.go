package model

type DbChanges struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"internal_id"`
	WebUserID      int    `gorm:"not null;default:0" json:"web_user_id"`
	ModelName      string `gorm:"not null" json:"model_name"`
	DataFrom       string `gorm:"not null" json:"data_from"`
	DataTo         string `gorm:"not null" json:"data_to"`
	AddedTimestamp int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}
