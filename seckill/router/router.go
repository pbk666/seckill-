package router

import (
	"decleration/controller"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRouter(h *server.Hertz) {
	seckillGroup := h.Group("/seckill")
	{
		seckillGroup.POST("", controller.SeckillHandler)
		seckillGroup.POST("/product/create", controller.CreateProductHandler)
		seckillGroup.POST("/order/create", controller.CreateOrderHandler)
	}
}
