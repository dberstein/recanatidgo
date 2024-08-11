package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// validResponse writes access token or error
func validResponse(c *gin.Context, user *UserCredentials) {
	token, err := CreateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetPwHash returns pwhash stored in database for username
func GetPwhash(username string) (string, error) {
	var pwhash string

	row := db.QueryRow(`SELECT pwhash FROM users WHERE username=?`, username)
	err := row.Scan(&pwhash)
	if err != nil {
		return pwhash, err
	}

	return pwhash, nil
}

func loginHandler(c *gin.Context) {
	user := &UserCredentials{}
	if err := c.BindJSON(user); err != nil { // Unmarshall request body ...
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pwhash, err := GetPwhash(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	valid := CheckPasswordHash(user.Password, pwhash)
	if valid {
		validResponse(c, user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
