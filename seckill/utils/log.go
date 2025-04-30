package utils

import (
	"context"
	"decleration/dao"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	// 配置日志输出到文件（自动创建）
	logger := &lumberjack.Logger{
		Filename:   "./seckill.log", // 日志文件路径
		MaxSize:    100,             // 单文件最大大小（MB）
		MaxBackups: 10,              // 保留旧日志文件的最大数量
		MaxAge:     30,              // 保留旧日志文件的最大天数
		Compress:   true,            // 是否压缩旧日志
	}

	hlog.SetOutput(logger) // 替换默认输出（控制台）为文件
}

func SeckillLogger(ctx context.Context, productID int, userID int) {
	db := dao.InitDB1()
	productDAO := dao.NewProductDAO(db)
	product, err := productDAO.GetProductByID(ctx, productID)
	if err != nil {
		return
	}
	orderDAO := dao.NewOrderDAO(db)
	order, err := orderDAO.GetOrdersByUserID(ctx, userID)
	hlog.Infof("user%d seckill successfully,the order is:%v,seckill product is:%v", userID, order, product)
	hlog.Infof("product StockSurplus decrease:1,leave:%d", product.StockSurplus)
}
