package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCurrentPage(c *gin.Context) (int, error) {
	p, err := strconv.Atoi(c.DefaultQuery("p", "1"))
	if err != nil {
		return 1, err
	}
	if p < 1 {
		p = 1
	}
	return p, nil
}

func DataHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")

		p, err := getCurrentPage(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ps := 10
		skip := (p - 1) * ps

		owm := NewOwm(os.Getenv("OWM_API_KEY"))
		data, err := owm.Query(c.Query("location"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %v! This is a protected route.", username),
			"page":    p,
			"limit":   fmt.Sprintf("%d,%d", skip, ps+1),
			"owm":     data,
		})
	}
}
