package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	stackmanager2 "github.com/h8r-dev/heighliner/pkg/stack/stackmanager"
)

type StackHandler interface {
	GetStack(c *gin.Context)
	ListStack(c *gin.Context)
}

func NewStackHandler(sm stackmanager2.StackManager) StackHandler {
	return &stackHandler{sm: sm}
}

type stackHandler struct {
	sm stackmanager2.StackManager
}

func (h *stackHandler) GetStack(c *gin.Context) {
	id := c.Param("id")
	s := h.sm.GetStack(id)
	if s == nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "stack id not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stack": s})
}

func (h *stackHandler) ListStack(c *gin.Context) {
	names := h.sm.ListStackNames()
	c.JSON(http.StatusOK, gin.H{"stacks": names})
}
