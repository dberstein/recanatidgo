package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// validLoginResponse writes access token or error
func validLoginResponse(c *gin.Context, user *UserCredentials) {
	token, err := createToken(user.Username)
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

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &UserCredentials{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pwhash, err := getPwhash(db, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		valid := checkPasswordHash(user.Password, pwhash)
		if valid {
			validLoginResponse(c, user)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
	}
}
