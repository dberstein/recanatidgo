package mw

import (
	"github.com/dberstein/recanatid-go/svc/token"
	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"
)

// rateLimitByToken is a middleware that rate limits according to TokenBucket and from request's JWT token
func RateLimitByTokenMiddleware(tb *ginratelimit.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := token.GetBearerToken(c.GetHeader("Authorization"))
		if token != "" {
			rate := ginratelimit.RateLimitByUserId(tb, token)
			rate(c)
		}
	}
}

func RateLimitByUserMiddleware(tb *ginratelimit.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if exists && username != "" {
			rate := ginratelimit.RateLimitByUserId(tb, username.(string))
			rate(c)
		}
	}
}
