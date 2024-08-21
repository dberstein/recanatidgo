package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func DataHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")

		owm := NewOwm(os.Getenv("OWM_API_KEY"))
		data, err := owm.Query(c.Query("location"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %v! This is a protected route.", username),
			"owm":     data,
		})
	}
}
