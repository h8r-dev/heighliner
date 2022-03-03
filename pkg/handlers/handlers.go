package handlers

import (
	"github.com/gin-gonic/gin"
)

// Health ...
func Health(c *gin.Context) {
	c.JSON(200, gin.H{})
}
