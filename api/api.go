package api

import (
	"github.com/chuccp/v2rayAuto/v2ray"
	"github.com/gin-gonic/gin"
)

func Subscribe(c *gin.Context) {
	urls := make([]string, 0)
	v2ray.RangeServer(func(server *v2ray.Server) {
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
