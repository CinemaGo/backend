package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Main Page",
	})
}
