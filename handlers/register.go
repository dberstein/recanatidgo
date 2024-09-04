package handlers

import (
	"database/sql"
	"net/http"

	"github.com/dberstein/recanatid-go/models"
	"github.com/dberstein/recanatid-go/svc/token"
	"github.com/dberstein/recanatid-go/typ"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(db *sql.DB, jwtMaker token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &typ.RegisterUser{}
		if err := c.BindJSON(user); err != nil { // Unmarshall request body ...
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		profile := models.NewProfile(db)
		if err := profile.Validate(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hasher := models.NewHasher()
		pwhash, err := hasher.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = profile.Insert(user, pwhash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		validLoginResponse(c, &typ.UserCredentials{Username: user.Username, Role: user.Role}, jwtMaker)
	}
}
