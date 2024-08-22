package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &RegisterUser{}
		if err := c.BindJSON(user); err != nil { // Unmarshall request body ...
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: username"})
		}
		if user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: password"})
		}
		if user.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: email"})
		}
		// if user.Role == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: role"})
		// }

		pwhash, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = db.Exec(
			"INSERT INTO users (username, pwhash, email, role) VALUES (?, ?, ?, ?)",
			&user.Username, &pwhash, &user.Email, &user.Role,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		validLoginResponse(c, &UserCredentials{Username: user.Username})
	}
}
