package controller

import (
	"context"
	"decleration/service"
	"decleration/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
)

func SeckillHandler(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ProductID int `json:"product_id"`
		UserID    int `json:"user_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		utils.Error(c, 400, "参数绑定失败", err)
	}

	utils.SeckillLogger(ctx, req.ProductID, req.UserID)

	product, err := service.GetProductByID(ctx, req.ProductID)
	if err != nil {
		log.Println("get product failed:", err)
		utils.Error(c, 500, "获取商品失败", err)
		return
	}
	if product == nil {
		utils.Error(c, 404, "商品不存在", nil)
		return
	}

	if err := service.SnapUpGoods(ctx, req.UserID, req.ProductID); err != nil {
		utils.Error(c, 500, "秒杀失败", err)
		return
	}

	utils.HResponse(c, 200, "秒杀成功", nil)
}
