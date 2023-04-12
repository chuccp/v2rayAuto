package main

import (
	"github.com/chuccp/v2rayAuto/api"
	"github.com/chuccp/v2rayAuto/v2ray"
	"github.com/gin-gonic/gin"
	"github.com/v2fly/v2ray-core/v5/common"
)

func main() {
	wsc := v2ray.CreateWebSocketConfig("127.0.0.1", 8400, 8800, 8, "WebSocketOverTls")
	server, err := v2ray.CreateWebSocketServer(wsc)
	common.Must(err)
	err = server.Start()
	common.Must(err)
	v2ray.RegisterServer(server)
	r := gin.Default()
	r.GET("/d3MuY2oyMDIw.md", api.Subscribe)
	r.Run(":8082")
}
