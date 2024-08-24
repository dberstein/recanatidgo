package handler

import (
	"database/sql"
	"net/http"

	"github.com/dberstein/recanatid-go/src/hash"
	"github.com/dberstein/recanatid-go/src/token"
	"github.com/dberstein/recanatid-go/src/typ"
	"github.com/gin-gonic/gin"
)

// validLoginResponse writes access token or error
func validLoginResponse(c *gin.Context, user *typ.UserCredentials, jwtMaker *token.JWTMaker) {
	token, err := jwtMaker.CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// getPwHash returns pwhash stored in database for username
func getPwhash(db *sql.DB, username string) (string, error) {
	var pwhash string

	row := db.QueryRow(`SELECT pwhash FROM users WHERE username=?`, username)
	err := row.Scan(&pwhash)
	if err != nil {
		return pwhash, err
	}

	return pwhash, nil
}

func LoginHandler(db *sql.DB, jwtMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &typ.UserCredentials{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pwhash, err := getPwhash(db, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		valid := hash.CheckPasswordHash(user.Password, pwhash)
		if valid {
			validLoginResponse(c, user, jwtMaker)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	}
}
