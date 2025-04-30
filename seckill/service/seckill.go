package service

import (
	"context"
	"decleration/dao"
	"decleration/kafka"
	"decleration/lock"
	"decleration/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func SnapUpGoods(ctx context.Context, userID int, productID int) error {
	db := dao.InitDB1()
	tx := db.Begin()
	productDAO := dao.NewProductDAO(db)

	// 查商品
	product, err := productDAO.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("查找商品失败: %v", err)
	}
	if product == nil {
		return errors.New("商品不存在")
	}

	mutex, err := lock.AcquireLock(fmt.Sprintf("seckill:lock:%d", product.ID), 10*time.Second)
	if err != nil {
		return err
	}
	defer lock.ReleaseLock(mutex)

	now := time.Now()
	if now.Before(product.StartTime) {
		return errors.New("秒杀尚未开始")
	}
	if now.After(product.EndTime) {
		return errors.New("秒杀已结束")
	}

	if err := productDAO.DecreaseStockSurplus(ctx, product.ID, 1); err != nil {
		tx.Rollback()
		return fmt.Errorf("扣减库存失败: %v", err)
	}

	orderMsg := model.OrderMessage{
		UserID:    userID,
		ProductID: productID,
		Quantity:  1,
	}
	msgBytes, _ := json.Marshal(orderMsg)
	if err := kafka.SendMessage("seckill-order", msgBytes); err != nil {
		tx.Rollback()
		return fmt.Errorf("发送 Kafka 消息失败: %v", err)
	}

	return nil
}
