package handler

import (
	"database/sql"
	"net/http"

	"github.com/dberstein/recanatid-go/hash"
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
		if user.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: username"})
			return
		}
		if user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: password"})
			return
		}
		if user.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: email"})
			return
		}
		// if user.Role == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing: role"})
		//  return
		// }

		pwhash, err := hash.HashPassword(user.Password)
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

		validLoginResponse(c, &typ.UserCredentials{Username: user.Username}, jwtMaker)
	}
}
