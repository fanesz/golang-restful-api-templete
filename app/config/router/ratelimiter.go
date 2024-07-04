package router

import (
	"backend/app/pkg/handler"
	"net/http"
	"time"

	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

const REQUEST_LIMIT = 10
const LIMIT_INTERVAL = time.Second * 5

func rateLimiterConfig() gin.HandlerFunc {
	return ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP()
		},
		LimitedHandler: func(c *gin.Context) {
			handler.Error(c, http.StatusTooManyRequests, "Too many requests")
			c.Abort()
		},
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return LIMIT_INTERVAL, REQUEST_LIMIT
		},
	})
}
