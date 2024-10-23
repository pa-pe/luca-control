package model

type TgUser struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName        string `gorm:"type:varchar(50);not null" json:"user_name"`
	FirstName       string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName        string `gorm:"type:varchar(100);not null" json:"last_name"`
	LanguageCode    string `gorm:"type:char(2);not null" json:"language_code"`
	IsBot           bool   `gorm:"not null;default:0" json:"is_bot,omitempty"`
	ChatbotPermit   byte   `gorm:"not null;default:0" json:"chatbot_permit"`
	TgCbFlowStepId  int    `gorm:"not null" json:"tg_cb_flow_step_id"`
	SrvsEmployeesId int    `gorm:"not null" json:"srvs_employees_id"`
	ShiftState      byte   `gorm:"not null;default:0" json:"shift_state"`
	AddedTimestamp  int64  `gorm:"autoCreateTime" json:"added_timestamp"`
}

//func (TgUser) TableName() string {
//	return "tg_users"
//}
