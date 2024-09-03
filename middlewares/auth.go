package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dberstein/recanatid-go/svc/token"

	"github.com/dgrijalva/jwt-go"
)

func isRoleNotAllowed(hasRole, requiredRole string) bool {
	return !(requiredRole == "" || hasRole == requiredRole)
}

func AuthMiddleware(jwtMaker *token.JWTMaker, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := token.GetBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token, err := jwtMaker.Parse(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if isRoleNotAllowed(claims["role"].(string), role) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Required token claims role: '%s' (has '%s')", role, claims["role"])})
				c.Abort()
			}
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}
