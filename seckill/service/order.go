package service

import (
	"context"
	"decleration/dao"
	"decleration/model"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateOrder(ctx context.Context, msg model.OrderMessage) error {
	db := dao.InitDB1()
	orderID := uuid.NewString()

	var product model.Product
	if err := db.WithContext(ctx).First(&product, msg.ProductID).Error; err != nil {
		return fmt.Errorf("找不到商品: %v", err)
	}

	if product.StockSurplus < msg.Quantity {
		return fmt.Errorf("库存不足")
	}

	orderItem := model.OrderItem{
		OrderID:   orderID,
		ProductID: msg.ProductID,
		Quantity:  msg.Quantity,
	}

	order := model.Order{
		OrderID: orderID,
		UserID:  msg.UserID,
		Orders:  []model.OrderItem{orderItem},
		Address: "默认地址",
		Total:   product.Price * float64(msg.Quantity),
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新库存
		product.StockSurplus -= msg.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("更新库存失败: %v", err)
		}

		// 创建订单
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败: %v", err)
		}

		return nil
	})
}
