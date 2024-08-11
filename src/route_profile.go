package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getProfileHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if exists {
		c.JSON(http.StatusOK, gin.H{"username": username.(string)})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
	}
}

func putProfileHandler(c *gin.Context) {
}
