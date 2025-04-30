package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
)

type H map[string]any

func Error(c *app.RequestContext, code int, msg string, err error) {
	c.JSON(code, H{
		"code": code,
		"msg":  msg,
		"err":  err.Error(),
	})
}

func HResponse(c *app.RequestContext, code int, msg string, data any) {
	c.JSON(code, H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
