package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/h8r-dev/heighliner/pkg/forms"
	stackmanager2 "github.com/h8r-dev/heighliner/pkg/stack/stackmanager"
)

// ApplicationHandler ...
type ApplicationHandler interface {
	PostApplication(c *gin.Context)
	GetApplication(c *gin.Context)
	ListApplication(c *gin.Context)
}

// NewApplicationHandler ...
func NewApplicationHandler(sm stackmanager2.StackManager) ApplicationHandler {
	return &applicationHandler{
		sm: sm,
	}
}

type applicationHandler struct {
	sm stackmanager2.StackManager
}

func (h *applicationHandler) PostApplication(c *gin.Context) {
	id := c.Param("id")
	var form forms.ApplicationForm
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	app, err := h.sm.InstantiateStack(id, form.StackID, form.Parameters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"application": app})
}

func (h *applicationHandler) GetApplication(c *gin.Context) {
	panic("implement me")
}

func (h *applicationHandler) ListApplication(c *gin.Context) {
	panic("implement me")
}
