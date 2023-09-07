package api

import (
	"bytes"
	"encoding/base64"
	"github.com/chuccp/v2rayAuto/core"
	"github.com/gin-gonic/gin"
	"github.com/v2fly/v2ray-core/v5/common"
	"net/http"
	"strconv"
)

type Server struct {
	context    *core.Context
	httpServer *http.Server
}

func (s *Server) Subscribe(c *gin.Context) {
	buff := new(bytes.Buffer)
	urls := make([]string, 0)
	s.context.RangeServer(func(server core.Server) {
		client := server.GetClient()
		for _, url := range client {
			urls = append(urls, url)
		}
	})
	lSize := len(urls) - 1
	for i, url := range urls {
		buff.WriteString(url)
		if i != lSize {
			buff.WriteString("\r\n")
		}
	}
	c.Writer.WriteString(base64.StdEncoding.EncodeToString(buff.Bytes()))
}
func (s *Server) Flush(c *gin.Context) {
	s.context.RangeServer(func(server core.Server) {
		server.Flush()
	})
}
func (s *Server) Start(context *core.Context) error {
	s.context = context
	if s.httpServer != nil {
		err := s.httpServer.Close()
		if err != nil {
			return err
		}
	}
	subscribe := common.Must2(s.context.ReadString("core", "subscribe")).(string)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET(subscribe, s.Subscribe)
	r.GET(subscribe+"_flush", s.Flush)
	cer := context.GetCertificate()
	s.httpServer = &http.Server{Addr: ":" + strconv.Itoa(s.context.GetPort()), Handler: r}
	return s.httpServer.ListenAndServeTLS(cer.CertificateFile, cer.KeyFile)
}
