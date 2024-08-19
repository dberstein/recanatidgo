package main

import (
	"github.com/gin-gonic/gin"
	ginratelimit "github.com/ljahier/gin-ratelimit"
)

// rateLimitByToken is a middleware that rate limits according to TokenBucket and from request's JWT token
func rateLimitByTokenMiddleware(tb *ginratelimit.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getBearerToken(c.GetHeader("Authorization"))
		if token != "" {
			rate := ginratelimit.RateLimitByUserId(tb, token)
			rate(c)
		}
	}
}

func rateLimitByUserMiddleware(tb *ginratelimit.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username != "" {
			rate := ginratelimit.RateLimitByUserId(tb, username.(string))
			rate(c)
		}
	}
}
