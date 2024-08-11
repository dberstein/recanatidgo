package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DataHandler(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %v! This is a protected route.", username)})
}
