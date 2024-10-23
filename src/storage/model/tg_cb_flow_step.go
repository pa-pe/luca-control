package model

type TgCbFlowStep struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	TgCbFlowId     int    `gorm:"not null;default:0" json:"tg_cb_flow_id"`
	Msg            string `gorm:"not null" json:"msg"`
	HandlerName    string `gorm:"not null" json:"handler_name"`
	RowOrder       int    `gorm:"not null;default:0" json:"row_order"`
	AddedTimestamp int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}
