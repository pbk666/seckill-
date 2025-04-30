package controller

import (
	"context"
	"decleration/model"
	"decleration/service"
	"decleration/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

func CreateProductHandler(ctx context.Context, c *app.RequestContext) {
	var product model.Product
	if err := c.BindAndValidate(&product); err != nil {
		utils.Error(c, 400, "参数绑定失败", err)
		return
	}

	if err := service.CreatProduct(product, ctx); err != nil {
		utils.Error(c, 500, "创建商品失败", err)
		return
	}

	utils.HResponse(c, 200, "创建商品成功", nil)
}
