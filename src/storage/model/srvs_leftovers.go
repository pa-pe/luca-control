package model

type SrvsLeftovers struct {
	ID                 int   `gorm:"primaryKey;autoIncrement" json:"id"`
	SrvsLocationId     int   `gorm:"not null;default:0" json:"srvs_location_id"`
	Date               int64 `gorm:"autoCreateTime" json:"date"`
	SrvsGoodsId        int   `gorm:"not null;default:0" json:"srvs_goods_id"`
	SrvsEmployeesId    int   `gorm:"not null;default:0" json:"srvs_employees_id"`
	QuantityStart      int   `gorm:"not null;default:0" json:"quantity_start"`
	QuantityEnd        int   `gorm:"not null;default:0" json:"quantity_end"`
	QuantitySell       int   `gorm:"not null;default:0" json:"quantity_sell"`
	QuantityWrittenOff int   `gorm:"not null;default:0" json:"quantity_written_off"`
}
