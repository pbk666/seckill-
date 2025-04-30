package dao

import (
	"decleration/model"
	"log"
)

func AutoMigrate() {
	if DB == nil {
		log.Fatal("数据库未初始化")
	}
	err := DB.AutoMigrate(
		&model.Product{},
		&model.Order{},
		&model.OrderMessage{},
		&model.OrderItem{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate 失败: %v", err)
	}
}
