package api

import (
	"github.com/chuccp/v2rayAuto/core"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Server struct {
}

func (s *Server) Subscribe(c *gin.Context) {

}
func (s *Server) Start(context *core.Context) error {
	r := gin.Default()
	r.GET("/d3MuY2oyMDIw.md", s.Subscribe)
	return r.Run(":" + strconv.Itoa(int(8000)))
}
