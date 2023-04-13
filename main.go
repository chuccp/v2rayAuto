package main

import (
	"github.com/chuccp/v2rayAuto/api"
	"github.com/chuccp/v2rayAuto/config"
	"github.com/chuccp/v2rayAuto/core"
	"github.com/chuccp/v2rayAuto/v2ray"
	"github.com/v2fly/v2ray-core/v5/common"
)

func main() {

	//wsc := v2ray.CreateWebSocketConfig("127.0.0.1", "cooge123@gmail.com", 8080, 8090, 8, "WebSocketOverTls")
	//server, err := v2ray.CreateWebSocketServer(wsc)
	//common.Must(err)
	//err = server.Start()
	//common.Must(err)
	//core.RegisterServer(server)
	//r := gin.Default()
	//r.GET("/d3MuY2oyMDIw.md", api.Subscribe)
	//r.Run(":" + strconv.Itoa(int(8000)))

	readConfig, err := config.ReadConfig("config.ini")
	common.Must(err)
	v2 := core.LoadConfig(readConfig, &api.Server{})
	v2.RegisterServer(&v2ray.WsServer{})
	err = v2.Start()
	common.Must(err)
}
