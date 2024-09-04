package handlers

import (
	"database/sql"
	"net/http"

	"github.com/dberstein/recanatid-go/models"
	"github.com/dberstein/recanatid-go/svc/token"
	"github.com/dberstein/recanatid-go/typ"
	"github.com/gin-gonic/gin"
)

// validLoginResponse writes access token or error
func validLoginResponse(c *gin.Context, user *typ.UserCredentials, jwtMaker token.JWTMaker) {
	token, err := jwtMaker.CreateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(db *sql.DB, jwtMaker token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &typ.UserCredentials{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hasher := models.NewHasher()
		pwhash, role, err := hasher.GetPwhashRole(db, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Role = role
		if hasher.CheckPasswordHash(user.Password, pwhash) {
			validLoginResponse(c, user, jwtMaker)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	}
}
