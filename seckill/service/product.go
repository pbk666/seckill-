package service

import (
	"context"
	"decleration/dao"
	"decleration/model"
	"fmt"
)

// 创建商品
func CreatProduct(product model.Product, ctx context.Context) error {
	//连接数据库
	db := dao.InitDB1()
	productDAO := dao.NewProductDAO(db)

	//创建商品
	err := productDAO.CreateProduct(ctx, product)
	if err != nil {
		return fmt.Errorf("创建商品失败:%v", err)
	}

	return nil
}

// 获取商品信息
func GetProductByID(ctx context.Context, productID int) (*model.Product, error) {
	//连接数据库
	db := dao.InitDB1()
	productDAO := dao.NewProductDAO(db)

	// 查询秒杀商品信息（假设从数据库或缓存查）
	product, err := productDAO.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("查询失败:%v", err)
	}
	return product, nil
}
