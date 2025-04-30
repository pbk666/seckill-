package dao

import (
	"context"
	"decleration/model"
	"fmt"
	"gorm.io/gorm"
)

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *OrderDAO {
	return &OrderDAO{db: db}
}

// 根据用户ID获取所有订单
func (dao *OrderDAO) GetOrdersByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	var orders []model.Order
	err := dao.db.WithContext(ctx).
		Preload("Orders").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

// 创建订单（包含订单项）事务处理
func (dao *OrderDAO) CreateOrderWithItems(ctx context.Context, order *model.Order) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("订单创建失败: %v", err)
		}
		return nil
	})
}

func (dao *OrderDAO) CreateOrder(ctx context.Context, order *model.Order) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败: %v", err)
		}
		return nil
	})
}
