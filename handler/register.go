package handler

import (
	"database/sql"
	"net/http"

	"github.com/dberstein/recanatid-go/hash"
	"github.com/dberstein/recanatid-go/model"
	"github.com/dberstein/recanatid-go/token"
	"github.com/dberstein/recanatid-go/typ"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(db *sql.DB, jwtMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &typ.RegisterUser{}
		if err := c.BindJSON(user); err != nil { // Unmarshall request body ...
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := model.ValidateRegisterUser(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pwhash, err := hash.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = model.InsertRegisterUser(db, user, pwhash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		validLoginResponse(c, &typ.UserCredentials{Username: user.Username}, jwtMaker)
	}
}
