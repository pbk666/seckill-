package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

func LogMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 记录请求信息
		hlog.Infof("Received request: %s %s", ctx.Request.Method(), ctx.Request.RequestURI())

		// 记录处理开始时间
		start := time.Now()

		// 继续处理请求
		ctx.Next(c)

		// 记录返回状态码
		hlog.Infof("Received status code: %d", ctx.Response.StatusCode())

		// 记录处理时间
		duration := time.Since(start)
		hlog.Infof("System processing time: %v", duration)

		// 如果需要记录额外的系统处理信息，比如是否发生错误、外部服务调用等
		if ctx.Response.StatusCode() >= 400 {
			hlog.Errorf("Request failed with status code: %d", ctx.Response.StatusCode())
		} else {
			hlog.Infof("Request processed successfully")
		}
	}
}
