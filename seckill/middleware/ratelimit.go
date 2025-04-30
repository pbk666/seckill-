package middleware

import (
	"context"
	"decleration/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/juju/ratelimit"
	"time"
)

func RateLimitMiddleware(limit int64) app.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(time.Second, limit, limit)
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取当前的令牌桶令牌数量
		if bucket.TakeAvailable(1) == 1 {
			// 允许通过请求
			c.Next(ctx)
		} else {
			// 超过请求次数限制，返回限流错误
			utils.HResponse(c, 429, "message", "Too many requests")
		}
	}
}
