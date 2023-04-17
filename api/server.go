package api

import (
	"github.com/chuccp/v2rayAuto/core"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Server struct {
	context *core.Context
}

func (s *Server) Subscribe(c *gin.Context) {
	urls := make([]string, 0)
	s.context.RangeServer(func(server core.Server) {
		client := server.GetClient()
		for _, url := range client {
			urls = append(urls, url)
		}
	})
	for i, url := range urls {
		c.Writer.WriteString(url)
		if i != len(urls) {
			c.Writer.WriteString("\n")
		}
	}
}
func (s *Server) Flush(c *gin.Context) {
	s.context.RangeServer(func(server core.Server) {
		server.Flush()
	})
}
func (s *Server) Start(context *core.Context) error {
	s.context = context
	r := gin.Default()
	r.GET("/d3MuY2oyMDIw.md", s.Subscribe)
	r.GET("/d3MuY2oyMDIwFlush.md", s.Flush)
	readInt, err := context.ReadInt("web", "port")
	if err != nil {
		return err
	}
	return r.Run(":" + strconv.Itoa(readInt))
}
