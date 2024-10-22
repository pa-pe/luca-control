package model

type SrvsShifts struct {
	ID                int   `gorm:"primaryKey;autoIncrement" json:"id"`
	SrvsLocationId    int   `gorm:"not null;default:0" json:"srvs_location_id"`
	Date              int64 `gorm:"autoCreateTime" json:"date"`
	SrvsEmployeesId   int   `gorm:"not null;default:0" json:"srvs_employees_id"`
	Salary            int   `gorm:"not null;default:0" json:"salary"`
	Paid              int   `gorm:"not null;default:0" json:"paid"`
	LeftToPay         int   `gorm:"not null;default:0" json:"left_to_pay"`
	Tips              int   `gorm:"not null;default:0" json:"tips"`
	QuantityPostCards int   `gorm:"not null;default:0" json:"quantity_post_cards"`
	QuantityPrints    int   `gorm:"not null;default:0" json:"quantity_prints"`
	QuantityFeedbacks int   `gorm:"not null;default:0" json:"quantity_feedbacks"`
}
