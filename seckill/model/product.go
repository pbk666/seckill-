package model

import "time"

type Product struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:text"`
	Type        string  `json:"type" gorm:"type:varchar(50)"`
	Price       float64 `json:"price" gorm:"not null"`

	Stock        int `json:"stock" gorm:"not null"`
	StockSurplus int `json:"stock_surplus" gorm:"not null"`
	LimitPerUser int `json:"limit_per_user" gorm:"default:1"` // ✅ 每人限购

	Status int `json:"status" gorm:"default:1"` // 1:上架, 0:下架

	StartTime time.Time `json:"start_time"` // ✅ 秒杀开始时间（时间戳）
	EndTime   time.Time `json:"end_time"`   // ✅ 秒杀结束时间（时间戳）

	CreatedAt time.Time
	UpdatedAt time.Time
}
