package model

type SrvsLocationList struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string `gorm:"not null" json:"name"`
	RowOrder       int    `gorm:"not null;default:0" json:"row_order"`
	AddedTimestamp int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}
