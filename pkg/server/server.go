package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/h8r-dev/heighliner/pkg/handlers"
	stackmanager2 "github.com/h8r-dev/heighliner/pkg/stack/stackmanager"
)

// Server ...
type Server interface {
	Start() error
}

type server struct {
	port int
	sm   stackmanager2.StackManager
}

func (s *server) Start() error {
	r := gin.Default()
	s.registerV1APIs(r.Group("/v1"))

	addr := fmt.Sprintf(":%d", s.port)
	return r.Run(addr)
}

func (s *server) registerV1APIs(r *gin.RouterGroup) {
	r.GET("/health", handlers.Health)

	appGroup := r.Group("/applications")
	ah := handlers.NewApplicationHandler(s.sm)
	appGroup.POST("/:id", ah.PostApplication)
	appGroup.GET("/", ah.ListApplication)
	appGroup.GET("/:id/*action", ah.GetApplication)

	stackGroup := r.Group("/stacks")
	sh := handlers.NewStackHandler(s.sm)
	stackGroup.GET("/", sh.ListStack)
	stackGroup.GET("/:id", sh.GetStack)
}

// New ...
func New(port int) Server {
	s := &server{
		port: port,
		sm:   stackmanager2.New(),
	}
	return s
}
