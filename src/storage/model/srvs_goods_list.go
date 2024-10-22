package model

type SrvsGoodsList struct {
	ID             int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string  `gorm:"not null" json:"name"`
	Price          float32 `gorm:"not null" json:"price"`
	AddedTimestamp int64   `gorm:"autoCreateTime" json:"added_timestamp"`
}
