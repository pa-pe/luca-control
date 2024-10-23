package model

type SrvsEmployeesList struct {
	ID             int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string  `gorm:"not null" json:"name"`
	Percentage     float32 `gorm:"not null" json:"percentage"`
	AddedTimestamp int64   `gorm:"autoCreateTime" json:"added_timestamp"`
}

func (SrvsEmployeesList) TableName() string {
	return "srvs_employees_list"
}
