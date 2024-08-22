package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getProfileHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		var user RegisterUser
		row := db.QueryRow(`SELECT username, email, role FROM users WHERE username=?`, username)
		if err := row.Scan(&user.Username, &user.Email, &user.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func putProfileHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &RegisterUser{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		if user.Email != "" {
			_, err := db.Exec(`UPDATE users SET email = ? WHERE username = ?`, &user.Email, username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if user.Password != "" {
			pwhash, err := hashPassword(user.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_, err = db.Exec(`UPDATE users SET pwhash = ? WHERE username = ?`, &pwhash, username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if user.Role != "" {
			_, err := db.Exec(`UPDATE users SET role = ? WHERE username = ?`, &user.Role, username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
}
