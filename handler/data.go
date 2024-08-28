package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dberstein/recanatid-go/model"
	"github.com/dberstein/recanatid-go/svc/owm"
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

func DataHandler(db *sql.DB, o owm.Owmer) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")

		p, err := getCurrentPage(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		owmData, err := o.Query(c.Query("location"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := model.NewData(db, 3)
		persons, err := data.ListUsers(p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %v! This is a protected route.", username),
			"owm":     owmData,
			"page":    p,
			"persons": persons,
		})
	}
}
