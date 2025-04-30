package controller

import (
	"context"
	"decleration/model"
	"decleration/service"
	"decleration/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

func CreateOrderHandler(ctx context.Context, c *app.RequestContext) {
	var msg model.OrderMessage
	if err := c.BindAndValidate(&msg); err != nil {
		utils.Error(c, 400, "参数绑定失败", err)
	}
	if err := service.CreateOrder(ctx, msg); err != nil {
		utils.Error(c, 500, "创建订单失败", err)
	} else {
		utils.HResponse(c, 200, "订单创建成功", nil)
	}
}
