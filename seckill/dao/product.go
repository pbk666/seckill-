package dao

import (
	"context"
	"decleration/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProductDAO struct {
	db *gorm.DB
}

func NewProductDAO(db *gorm.DB) *ProductDAO {
	return &ProductDAO{db: db}
}

// 创建商品
func (dao *ProductDAO) CreateProduct(ctx context.Context, product model.Product) error {
	return dao.db.WithContext(ctx).Model(&model.Product{}).Create(&product).Error
}

// 查找商品信息（获取价格）
func (dao *ProductDAO) GetProductByID(ctx context.Context, productID int) (*model.Product, error) {
	var product model.Product
	if err := dao.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("查找商品失败: %v", err)
	}
	return &product, nil
}

// 扣减库存
func (dao *ProductDAO) DecreaseStockSurplus(ctx context.Context, id int, quantity int) error {
	// 原子性地扣库存，防止超卖
	result := dao.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock_surplus >= ?", id, quantity).
		UpdateColumn("stock_surplus", gorm.Expr("stock_surplus - ?", quantity))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("库存不足")
	}
	return nil

}

// 增加剩余库存（失败回滚时用）
func (dao *ProductDAO) IncreaseStockSurplus(ctx context.Context, productID int, quantity int) {
	dao.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock_surplus", gorm.Expr("stock_surplus + ?", quantity))
}

// 修改商品状态
func (dao *ProductDAO) UpdateProductStatus(ctx context.Context, id int, status int) error {
	return dao.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? ", id).
		Update("status", status).Error
}

// 更新秒杀时间（可用于后台配置秒杀活动）
func (dao *ProductDAO) UpdateSeckillTime(ctx context.Context, productID int, startTime, endTime int64) error {
	return dao.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
		}).Error
}
