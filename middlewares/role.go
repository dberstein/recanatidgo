package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var allow bool
		r := c.MustGet("role")
		for _, rr := range roles {
			if r == rr {
				allow = true
				break
			}
		}
		if !allow {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Required token claims role: %q has '%s')", roles, r)})
			c.Abort()
			return
		}

		c.Next()
	}
}
