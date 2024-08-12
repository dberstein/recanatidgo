package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	user := &RegisterUser{}
	if err := c.BindJSON(user); err != nil { // Unmarshall request body ...
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pwhash, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("sqlite3", "x.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(
		"INSERT INTO users (username, pwhash, email, role) VALUES (?, ?, ?, ?)",
		&user.Username, &pwhash, &user.Email, &user.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	validResponse(c, &UserCredentials{Username: user.Username})
}
